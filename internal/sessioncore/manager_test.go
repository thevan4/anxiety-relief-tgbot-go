package sessioncore

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestManager_Start(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx := m.Start(parentCtx, 123)
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
	if ctx.Err() != nil {
		t.Error("new context should not be cancelled")
	}
	if !m.Has(123) {
		t.Error("manager should track user 123")
	}
}

func TestManager_StartCancelsPrevious(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	// Start first session
	ctx1 := m.Start(parentCtx, 123)

	// Start second session for same user
	ctx2 := m.Start(parentCtx, 123)

	// First context should be cancelled
	if ctx1.Err() == nil {
		t.Error("first context should be cancelled")
	}

	// Second context should be active
	if ctx2.Err() != nil {
		t.Error("second context should not be cancelled")
	}
}

func TestManager_Cancel(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx := m.Start(parentCtx, 123)
	m.Cancel(123)

	if ctx.Err() == nil {
		t.Error("context should be cancelled after Cancel()")
	}
	if m.Has(123) {
		t.Error("user should be removed after Cancel()")
	}
}

func TestManager_CancelNonExistent(t *testing.T) {
	t.Parallel()
	m := NewManager()
	// Should not panic
	m.Cancel(999)
}

func TestManager_MultipleUsers(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx1 := m.Start(parentCtx, 1)
	ctx2 := m.Start(parentCtx, 2)
	ctx3 := m.Start(parentCtx, 3)

	if m.Count() != 3 {
		t.Errorf("expected 3 users, got %d", m.Count())
	}

	// Cancel one user
	m.Cancel(2)

	if ctx1.Err() != nil {
		t.Error("user 1 context should not be cancelled")
	}
	if ctx2.Err() == nil {
		t.Error("user 2 context should be cancelled")
	}
	if ctx3.Err() != nil {
		t.Error("user 3 context should not be cancelled")
	}
	if m.Count() != 2 {
		t.Errorf("expected 2 users after cancel, got %d", m.Count())
	}
}

func TestManager_CancelAll(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx1 := m.Start(parentCtx, 1)
	ctx2 := m.Start(parentCtx, 2)
	ctx3 := m.Start(parentCtx, 3)

	m.CancelAll()

	if ctx1.Err() == nil || ctx2.Err() == nil || ctx3.Err() == nil {
		t.Error("all contexts should be cancelled")
	}
	if m.Count() != 0 {
		t.Errorf("expected 0 users after CancelAll, got %d", m.Count())
	}
}

func TestManager_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(userID int64) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				m.Start(parentCtx, userID)
				m.Has(userID)
				time.Sleep(time.Microsecond)
			}
			m.Cancel(userID)
		}(int64(i))
	}
	wg.Wait()

	// Should not panic and should end with 0 users
	if m.Count() != 0 {
		t.Errorf("expected 0 users after concurrent test, got %d", m.Count())
	}
}

func TestManager_ContextInheritsCancellation(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx, parentCancel := context.WithCancel(context.Background())

	ctx := m.Start(parentCtx, 123)

	// Cancel parent
	parentCancel()

	// Child context should also be cancelled
	select {
	case <-ctx.Done():
		// OK
	case <-time.After(100 * time.Millisecond):
		t.Error("context should be cancelled when parent is cancelled")
	}
}

func TestManager_Has(t *testing.T) {
	t.Parallel()
	m := NewManager()

	if m.Has(123) {
		t.Error("should not have user 123 before start")
	}

	m.Start(context.Background(), 123)

	if !m.Has(123) {
		t.Error("should have user 123 after start")
	}

	m.Cancel(123)

	if m.Has(123) {
		t.Error("should not have user 123 after cancel")
	}
}
