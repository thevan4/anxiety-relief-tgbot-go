package cleanup

import (
	"context"
	"testing"
	"time"

	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
)

// mockStorage implements session.Storage for testing.
type mockStorage struct {
	states         map[int64]session.State
	menuIDs        map[int64]int
	menuCreatedAt  map[int64]time.Time
	cleanupQueue   map[int64]time.Time
	cleanupRetries map[int64]*session.CleanupRetry
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		states:         make(map[int64]session.State),
		menuIDs:        make(map[int64]int),
		menuCreatedAt:  make(map[int64]time.Time),
		cleanupQueue:   make(map[int64]time.Time),
		cleanupRetries: make(map[int64]*session.CleanupRetry),
	}
}

func (m *mockStorage) SetState(_ context.Context, userID int64, state session.State) error {
	m.states[userID] = state
	return nil
}

func (m *mockStorage) GetState(_ context.Context, userID int64) (session.State, error) {
	return m.states[userID], nil
}

func (m *mockStorage) ClearState(_ context.Context, userID int64) error {
	delete(m.states, userID)
	return nil
}

func (m *mockStorage) SetMessageID(_ context.Context, userID int64, messageID int) error {
	m.menuIDs[userID] = messageID
	return nil
}

func (m *mockStorage) GetMessageID(_ context.Context, userID int64) (int, error) {
	return m.menuIDs[userID], nil
}

func (m *mockStorage) SetHolderMessageID(_ context.Context, _ int64, _ int) error { return nil }
func (m *mockStorage) GetHolderMessageID(_ context.Context, _ int64) (int, error) { return 0, nil }
func (m *mockStorage) ClearHolderMessageID(_ context.Context, _ int64) error      { return nil }

func (m *mockStorage) SetMenuMessageID(_ context.Context, userID int64, messageID int) error {
	m.menuIDs[userID] = messageID
	return nil
}

func (m *mockStorage) GetMenuMessageID(_ context.Context, userID int64) (int, error) {
	return m.menuIDs[userID], nil
}

func (m *mockStorage) ClearMenuMessageID(_ context.Context, userID int64) error {
	delete(m.menuIDs, userID)
	return nil
}

func (m *mockStorage) SetMenuCreatedAt(_ context.Context, userID int64, t time.Time) error {
	m.menuCreatedAt[userID] = t
	return nil
}

func (m *mockStorage) GetMenuCreatedAt(_ context.Context, userID int64) (time.Time, error) {
	return m.menuCreatedAt[userID], nil
}

func (m *mockStorage) SetLang(_ context.Context, _ int64, _ string) error { return nil }
func (m *mockStorage) GetLang(_ context.Context, _ int64) (string, error) { return "en", nil }
func (m *mockStorage) ClearSession(_ context.Context, _ int64) error      { return nil }
func (m *mockStorage) Close() error                                       { return nil }

func (m *mockStorage) AddToCleanupQueue(_ context.Context, userID int64, checkAt time.Time) error {
	m.cleanupQueue[userID] = checkAt
	return nil
}

func (m *mockStorage) GetPendingCleanup(_ context.Context, now time.Time, limit int64) ([]int64, error) {
	var result []int64
	for userID, checkAt := range m.cleanupQueue {
		if checkAt.Before(now) || checkAt.Equal(now) {
			result = append(result, userID)
			if int64(len(result)) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *mockStorage) RemoveFromCleanupQueue(_ context.Context, userID int64) error {
	delete(m.cleanupQueue, userID)
	return nil
}

func (m *mockStorage) SetCleanupRetry(_ context.Context, userID int64, retry session.CleanupRetry) error {
	m.cleanupRetries[userID] = &retry
	return nil
}

func (m *mockStorage) GetCleanupRetry(_ context.Context, userID int64) (*session.CleanupRetry, error) {
	if r, ok := m.cleanupRetries[userID]; ok {
		return r, nil
	}
	return nil, session.ErrCleanupRetryNotFound
}

func (m *mockStorage) ClearCleanupRetry(_ context.Context, userID int64) error {
	delete(m.cleanupRetries, userID)
	return nil
}

func TestWorkerProcessesInactiveUser(t *testing.T) {
	t.Parallel()
	storage := newMockStorage()
	ctx := context.Background()

	// Setup: user with old menu, no state (inactive)
	userID := int64(123)
	storage.menuIDs[userID] = 456
	storage.menuCreatedAt[userID] = time.Now().Add(-47 * time.Hour)
	storage.cleanupQueue[userID] = time.Now().Add(-1 * time.Hour)

	// Create worker with mock time
	w := &Worker{
		ctx:     ctx,
		bot:     nil, // Will cause delete to fail, but that's ok for this test
		storage: storage,
		core:    session.NewCleanupProcessor(session.DefaultConfig()),
		timeNow: time.Now,
	}

	// Process - should decide to delete since user is inactive
	result := w.core.EvaluateCleanup(
		storage.menuCreatedAt[userID],
		string(storage.states[userID]),
		nil,
		time.Now(),
	)

	if result.Decision != session.DecisionDelete {
		t.Errorf("Expected DecisionDelete for inactive user, got %v", result.Decision)
	}
}

func TestWorkerSchedulesRetryForActiveUser(t *testing.T) {
	t.Parallel()
	storage := newMockStorage()

	// Setup: user with menu, active state
	userID := int64(123)
	storage.menuIDs[userID] = 456
	storage.menuCreatedAt[userID] = time.Now().Add(-47 * time.Hour)
	storage.states[userID] = session.StateBreathingRunning
	storage.cleanupQueue[userID] = time.Now().Add(-1 * time.Hour)

	processor := session.NewCleanupProcessor(session.DefaultConfig())

	// Process - should decide to retry since user is active
	result := processor.EvaluateCleanup(
		storage.menuCreatedAt[userID],
		string(storage.states[userID]),
		nil,
		time.Now(),
	)

	if result.Decision != session.DecisionRetry {
		t.Errorf("Expected DecisionRetry for active user, got %v", result.Decision)
	}
	if result.Retry == nil {
		t.Error("Expected Retry info to be set")
	}
}

func TestWorkerKeepsMenuWhenUserFinishes(t *testing.T) {
	t.Parallel()
	storage := newMockStorage()

	// Setup: user was active, now finished (state empty)
	userID := int64(123)
	storage.menuIDs[userID] = 456
	storage.menuCreatedAt[userID] = time.Now().Add(-47 * time.Hour)
	// state is empty (user finished)
	storage.cleanupRetries[userID] = &session.CleanupRetry{
		State:   string(session.StateBreathingRunning),
		Attempt: 1,
	}

	processor := session.NewCleanupProcessor(session.DefaultConfig())

	// Process - should decide to keep since user finished activity
	result := processor.EvaluateCleanup(
		storage.menuCreatedAt[userID],
		"", // empty = finished
		storage.cleanupRetries[userID],
		time.Now(),
	)

	if result.Decision != session.DecisionKeep {
		t.Errorf("Expected DecisionKeep when user finishes, got %v", result.Decision)
	}
}

func TestWorkerDeletesStuckUser(t *testing.T) {
	t.Parallel()
	storage := newMockStorage()

	// Setup: user stuck in same state
	userID := int64(123)
	storage.menuIDs[userID] = 456
	storage.menuCreatedAt[userID] = time.Now().Add(-47 * time.Hour)
	storage.states[userID] = session.StateBreathingRunning
	storage.cleanupRetries[userID] = &session.CleanupRetry{
		State:   string(session.StateBreathingRunning), // same state!
		Attempt: 1,
	}

	processor := session.NewCleanupProcessor(session.DefaultConfig())

	// Process - should delete since user is stuck
	result := processor.EvaluateCleanup(
		storage.menuCreatedAt[userID],
		string(storage.states[userID]),
		storage.cleanupRetries[userID],
		time.Now(),
	)

	if result.Decision != session.DecisionDelete {
		t.Errorf("Expected DecisionDelete for stuck user, got %v", result.Decision)
	}
}

func TestWorkerDeletesAfterMaxRetries(t *testing.T) {
	t.Parallel()
	storage := newMockStorage()

	// Setup: user keeps changing states, max retries reached
	userID := int64(123)
	storage.menuIDs[userID] = 456
	storage.menuCreatedAt[userID] = time.Now().Add(-47 * time.Hour)
	storage.states[userID] = session.StateGroundingStep3 // different from retry state
	storage.cleanupRetries[userID] = &session.CleanupRetry{
		State:   string(session.StateBreathingRunning),
		Attempt: 6, // max retries
	}

	processor := session.NewCleanupProcessor(session.DefaultConfig())

	// Process - should delete since max retries reached
	result := processor.EvaluateCleanup(
		storage.menuCreatedAt[userID],
		string(storage.states[userID]),
		storage.cleanupRetries[userID],
		time.Now(),
	)

	if result.Decision != session.DecisionDelete {
		t.Errorf("Expected DecisionDelete after max retries, got %v", result.Decision)
	}
}

func TestWorkerDeletesExpiredResource(t *testing.T) {
	t.Parallel()

	// Setup: resource older than 48h
	createdAt := time.Now().Add(-49 * time.Hour)

	processor := session.NewCleanupProcessor(session.DefaultConfig())

	// Process - should delete since resource is too old
	result := processor.EvaluateCleanup(
		createdAt,
		string(session.StateBreathingRunning), // even if active
		nil,
		time.Now(),
	)

	if result.Decision != session.DecisionDelete {
		t.Errorf("Expected DecisionDelete for expired resource, got %v", result.Decision)
	}
}
