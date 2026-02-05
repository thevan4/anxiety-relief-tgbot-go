// Package cleanup provides background message cleanup functionality.
package cleanup

import (
	"context"
	"log"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
)

const (
	// checkInterval — интервал проверки очереди.
	checkInterval = 5 * time.Minute

	// retryInterval — интервал между retry проверками.
	retryInterval = 5 * time.Minute

	// maxRetries — максимальное количество retry (6 × 5 мин = 30 мин).
	maxRetries = 6

	// batchSize — размер пакета для обработки.
	batchSize = 100
)

// Worker performs background cleanup of old menu messages.
type Worker struct {
	ctx     context.Context
	bot     *telego.Bot
	storage session.Storage
	timeNow func() time.Time
}

// NewWorker creates a new cleanup worker.
func NewWorker(ctx context.Context, bot *telego.Bot, storage session.Storage) *Worker {
	return &Worker{
		ctx:     ctx,
		bot:     bot,
		storage: storage,
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

// processUser handles cleanup for a single user.
func (w *Worker) processUser(userID int64, now time.Time) {
	// Get menu creation time
	menuCreatedAt, err := w.storage.GetMenuCreatedAt(w.ctx, userID)
	if err != nil {
		log.Printf("ERROR: cleanup get menu created: %v", err)
		w.cleanupUser(userID, "error getting menu created time")
		return
	}

	// If no menu creation time, clean up
	if menuCreatedAt.IsZero() {
		w.cleanupUser(userID, "no menu created time")
		return
	}

	// Step 1: Check if menu is too old (>48h)
	menuAge := now.Sub(menuCreatedAt)
	if menuAge > session.SessionMaxTTL {
		log.Printf("WARNING: menu too old for user %d (age: %v), cleaning up without delete", userID, menuAge)
		w.cleanupUser(userID, "menu too old")
		return
	}

	// Get current state
	state, err := w.storage.GetState(w.ctx, userID)
	if err != nil {
		log.Printf("ERROR: cleanup get state: %v", err)
		w.reschedule(userID, now, 1*time.Minute)
		return
	}

	// Get retry info
	retry, err := w.storage.GetCleanupRetry(w.ctx, userID)
	if err != nil {
		log.Printf("ERROR: cleanup get retry: %v", err)
		w.reschedule(userID, now, 1*time.Minute)
		return
	}

	// Step 2: First check (no retry exists)
	if retry == nil {
		w.handleFirstCheck(userID, state, now)
		return
	}

	// Step 3: Retry check
	w.handleRetryCheck(userID, state, retry, now)
}

// handleFirstCheck handles the first cleanup check for a user.
func (w *Worker) handleFirstCheck(userID int64, state session.State, now time.Time) {
	// User NOT in exercise (state empty)
	if state == session.StateUnknown || state == "" {
		w.deleteMenuAndCleanup(userID)
		return
	}

	// User IN exercise — schedule retry
	log.Printf("DEBUG: user %d in exercise %s, scheduling retry", userID, state)
	_ = w.storage.SetCleanupRetry(w.ctx, userID, session.CleanupRetry{
		State:   string(state),
		Attempt: 1,
	})
	w.reschedule(userID, now, retryInterval)
}

// handleRetryCheck handles retry cleanup checks.
func (w *Worker) handleRetryCheck(userID int64, state session.State, retry *session.CleanupRetry, now time.Time) {
	currentState := string(state)

	// User EXITED exercise (state empty)
	if state == session.StateUnknown || state == "" {
		// User was active, finished exercise — DON'T delete!
		log.Printf("DEBUG: user %d finished exercise, not deleting menu", userID)
		_ = w.storage.ClearCleanupRetry(w.ctx, userID)
		_ = w.storage.RemoveFromCleanupQueue(w.ctx, userID)
		return
	}

	// User in SAME exercise
	if currentState == retry.State {
		// Stuck or taking long — delete
		log.Printf("DEBUG: user %d stuck in %s, deleting menu", userID, currentState)
		w.deleteMenuAndCleanup(userID)
		return
	}

	// User in NEW exercise (state changed)
	if retry.Attempt < maxRetries {
		// Still active, reschedule
		log.Printf("DEBUG: user %d switched to %s, retry %d", userID, currentState, retry.Attempt+1)
		_ = w.storage.SetCleanupRetry(w.ctx, userID, session.CleanupRetry{
			State:   currentState,
			Attempt: retry.Attempt + 1,
		})
		w.reschedule(userID, now, retryInterval)
		return
	}

	// Max retries reached — force delete
	log.Printf("DEBUG: user %d max retries, force deleting menu", userID)
	w.deleteMenuAndCleanup(userID)
}

// deleteMenuAndCleanup deletes the menu message and cleans up Redis.
func (w *Worker) deleteMenuAndCleanup(userID int64) {
	menuID, err := w.storage.GetMenuMessageID(w.ctx, userID)
	if err == nil && menuID != 0 {
		if err := w.bot.DeleteMessage(w.ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(userID),
			MessageID: menuID,
		}); err != nil {
			// Log but continue cleanup — message might already be deleted or too old
			log.Printf("DEBUG: cleanup delete message %d for user %d: %v", menuID, userID, err)
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
