// Package cleanup provides background message cleanup functionality.
package cleanup

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
)

const (
	// checkInterval — интервал проверки очереди.
	checkInterval = 5 * time.Minute

	// batchSize — размер пакета для обработки.
	batchSize = 100
)

// Worker performs background cleanup of old menu messages.
type Worker struct {
	ctx     context.Context
	bot     *telego.Bot
	storage session.Storage
	core    *session.CleanupProcessor
	timeNow func() time.Time
}

// NewWorker creates a new cleanup worker.
func NewWorker(ctx context.Context, bot *telego.Bot, storage session.Storage) *Worker {
	return &Worker{
		ctx:     ctx,
		bot:     bot,
		storage: storage,
		core:    session.NewCleanupProcessor(session.DefaultConfig()),
		timeNow: time.Now,
	}
}

// Start runs the cleanup worker loop.
func (w *Worker) Start() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	log.Println("Cleanup worker started")

	for {
		select {
		case <-w.ctx.Done():
			log.Println("Cleanup worker stopped")
			return
		case <-ticker.C:
			w.processQueue()
		}
	}
}

// processQueue processes all users pending cleanup.
func (w *Worker) processQueue() {
	now := w.timeNow()

	userIDs, err := w.storage.GetPendingCleanup(w.ctx, now, batchSize)
	if err != nil {
		log.Printf("ERROR: cleanup get pending: %v", err)
		return
	}

	if len(userIDs) == 0 {
		return
	}

	log.Printf("DEBUG: cleanup processing %d users", len(userIDs))

	for _, userID := range userIDs {
		if w.ctx.Err() != nil {
			return
		}
		w.processUser(userID, now)
	}
}

// processUser handles cleanup for a single user using session cleanup logic.
func (w *Worker) processUser(userID int64, now time.Time) {
	// Get data needed for evaluation
	resourceCreatedAt, err := w.storage.GetMenuCreatedAt(w.ctx, userID)
	if err != nil {
		log.Printf("ERROR: cleanup get menu created: %v", err)
		w.executeDecision(userID, session.CleanupResult{
			Decision: session.DecisionDelete,
			Reason:   "error getting resource creation time",
		})
		return
	}

	state, err := w.storage.GetState(w.ctx, userID)
	if err != nil {
		log.Printf("ERROR: cleanup get state: %v", err)
		w.reschedule(userID, now, 1*time.Minute)
		return
	}
	currentState := string(state)

	retry, err := w.storage.GetCleanupRetry(w.ctx, userID)
	if err != nil {
		if errors.Is(err, session.ErrCleanupRetryNotFound) {
			retry = nil
		} else {
			log.Printf("ERROR: cleanup get retry: %v", err)
			w.reschedule(userID, now, 1*time.Minute)
			return
		}
	}

	// Use session cleanup processor to evaluate what to do
	result := w.core.EvaluateCleanup(resourceCreatedAt, currentState, retry, now)

	log.Printf("DEBUG: user %d cleanup decision: %v (%s)", userID, result.Decision, result.Reason)

	// Execute the decision
	w.executeDecision(userID, result)
}

// executeDecision performs the action based on cleanup evaluation result.
func (w *Worker) executeDecision(userID int64, result session.CleanupResult) {
	switch result.Decision {
	case session.DecisionDelete:
		w.deleteMenuAndCleanup(userID)

	case session.DecisionKeep:
		// User finished activity, keep the menu but remove from queue
		_ = w.storage.ClearCleanupRetry(w.ctx, userID)
		_ = w.storage.RemoveFromCleanupQueue(w.ctx, userID)

	case session.DecisionRetry:
		// Schedule retry check
		if result.Retry != nil {
			_ = w.storage.SetCleanupRetry(w.ctx, userID, *result.Retry)
		}
		_ = w.storage.AddToCleanupQueue(w.ctx, userID, result.NextCheck)
	}
}

// deleteMenuAndCleanup deletes the menu message and cleans up Redis.
func (w *Worker) deleteMenuAndCleanup(userID int64) {
	menuID, err := w.storage.GetMenuMessageID(w.ctx, userID)
	if err == nil && menuID != 0 {
		if delErr := w.bot.DeleteMessage(w.ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(userID),
			MessageID: menuID,
		}); delErr != nil {
			// Log but continue cleanup — message might already be deleted or too old
			log.Printf("DEBUG: cleanup delete message %d for user %d: %v", menuID, userID, delErr)
		}
	}

	w.cleanupUser(userID, "menu deleted")
}

// cleanupUser removes user from queue and clears session data.
func (w *Worker) cleanupUser(userID int64, reason string) {
	log.Printf("DEBUG: cleanup user %d: %s", userID, reason)

	_ = w.storage.RemoveFromCleanupQueue(w.ctx, userID)
	_ = w.storage.ClearCleanupRetry(w.ctx, userID)
	_ = w.storage.ClearState(w.ctx, userID)
	_ = w.storage.ClearMenuMessageID(w.ctx, userID)
}

// reschedule updates user's check time in the queue.
func (w *Worker) reschedule(userID int64, now time.Time, delay time.Duration) {
	nextCheck := now.Add(delay)
	if err := w.storage.AddToCleanupQueue(w.ctx, userID, nextCheck); err != nil {
		log.Printf("ERROR: cleanup reschedule user %d: %v", userID, err)
	}
}
