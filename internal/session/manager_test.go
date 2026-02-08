package session

import (
	"context"
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	if sm == nil {
		t.Fatal("NewManager returned nil")
	}
	if sm.cancels == nil {
		t.Error("cancels map is nil")
	}
}

func TestStartSession(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx := context.Background()

	ctx := sm.StartSession(parentCtx, 123)
	if ctx == nil {
		t.Fatal("StartSession returned nil context")
	}

	// Context should not be cancelled initially
	select {
	case <-ctx.Done():
		t.Error("new session context is already cancelled")
	default:
	}
}

func TestStartSessionCancelsPrevious(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx := context.Background()

	ctx1 := sm.StartSession(parentCtx, 123)
	ctx2 := sm.StartSession(parentCtx, 123)

	// First context should be cancelled
	select {
	case <-ctx1.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Error("previous session context was not cancelled")
	}

	// Second context should still be active
	select {
	case <-ctx2.Done():
		t.Error("new session context is cancelled")
	default:
	}
}

func TestCancelSession(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx := context.Background()

	ctx := sm.StartSession(parentCtx, 123)
	sm.CancelSession(123)

	// Context should be cancelled
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Error("session context was not cancelled")
	}
}

func TestCancelSessionNonExistent(t *testing.T) {
	t.Parallel()
	sm := NewManager()

	// Should not panic
	sm.CancelSession(999)
}

func TestGetSessionContext(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx := context.Background()

	// First call creates new session
	ctx1 := sm.GetSessionContext(parentCtx, 123)
	if ctx1 == nil {
		t.Fatal("GetSessionContext returned nil")
	}

	// Second call replaces session
	ctx2 := sm.GetSessionContext(parentCtx, 123)
	if ctx2 == nil {
		t.Fatal("GetSessionContext returned nil on second call")
	}
}

func TestMultipleUsers(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx := context.Background()

	ctx1 := sm.StartSession(parentCtx, 111)
	ctx2 := sm.StartSession(parentCtx, 222)
	ctx3 := sm.StartSession(parentCtx, 333)

	// Cancel user 222
	sm.CancelSession(222)

	// User 222 context should be cancelled
	select {
	case <-ctx2.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Error("user 222 context was not cancelled")
	}

	// Other contexts should still be active
	select {
	case <-ctx1.Done():
		t.Error("user 111 context was cancelled unexpectedly")
	default:
	}

	select {
	case <-ctx3.Done():
		t.Error("user 333 context was cancelled unexpectedly")
	default:
	}
}

func TestParentContextCancellation(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx, cancel := context.WithCancel(context.Background())

	ctx := sm.StartSession(parentCtx, 123)

	// Cancel parent context
	cancel()

	// Child context should be cancelled
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Error("child context was not cancelled when parent was cancelled")
	}
}

func TestConcurrentSessionOperations(t *testing.T) {
	t.Parallel()
	sm := NewManager()
	parentCtx := context.Background()
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(userID int64) {
			for j := 0; j < 100; j++ {
				sm.StartSession(parentCtx, userID)
				sm.GetSessionContext(parentCtx, userID)
				sm.CancelSession(userID)
			}
			done <- true
		}(int64(i))
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestNewManager(t *testing.T) {
	t.Parallel()
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if m.Count() != 0 {
		t.Error("new manager should have 0 sessions")
	}
}

func TestManagerHas(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	if m.Has(123) {
		t.Error("Has() should return false for non-existent user")
	}

	m.Start(parentCtx, 123)
	if !m.Has(123) {
		t.Error("Has() should return true after Start()")
	}

	m.Cancel(123)
	if m.Has(123) {
		t.Error("Has() should return false after Cancel()")
	}
}

func TestManagerCount(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	if m.Count() != 0 {
		t.Errorf("Count() = %d, want 0", m.Count())
	}

	m.Start(parentCtx, 1)
	m.Start(parentCtx, 2)
	m.Start(parentCtx, 3)

	if m.Count() != 3 {
		t.Errorf("Count() = %d, want 3", m.Count())
	}

	m.Cancel(2)
	if m.Count() != 2 {
		t.Errorf("Count() = %d, want 2 after cancel", m.Count())
	}
}

func TestManagerCancelAll(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx1 := m.Start(parentCtx, 1)
	ctx2 := m.Start(parentCtx, 2)
	ctx3 := m.Start(parentCtx, 3)

	m.CancelAll()

	// All contexts should be cancelled
	for i, ctx := range []context.Context{ctx1, ctx2, ctx3} {
		select {
		case <-ctx.Done():
			// expected
		case <-time.After(100 * time.Millisecond):
			t.Errorf("context %d was not cancelled by CancelAll()", i+1)
		}
	}

	if m.Count() != 0 {
		t.Errorf("Count() = %d after CancelAll(), want 0", m.Count())
	}
}

func TestManagerGetOrStart(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	// First call should create new session
	ctx1 := m.GetOrStart(parentCtx, 123)
	if ctx1 == nil {
		t.Fatal("GetOrStart returned nil")
	}
	if !m.Has(123) {
		t.Error("session should exist after GetOrStart")
	}

	// Second call should replace existing session
	ctx2 := m.GetOrStart(parentCtx, 123)
	if ctx2 == nil {
		t.Fatal("second GetOrStart returned nil")
	}

	// Still only one session for this user
	if m.Count() != 1 {
		t.Errorf("Count() = %d, want 1", m.Count())
	}
}

func TestManagerStart(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx := m.Start(parentCtx, 123)
	if ctx == nil {
		t.Fatal("Start returned nil")
	}

	// Context should be active
	select {
	case <-ctx.Done():
		t.Error("context should not be cancelled initially")
	default:
	}
}

func TestManagerCancel(t *testing.T) {
	t.Parallel()
	m := NewManager()
	parentCtx := context.Background()

	ctx := m.Start(parentCtx, 123)
	m.Cancel(123)

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Error("context was not cancelled")
	}
}

func TestManagerCancelNonExistent(t *testing.T) {
	t.Parallel()
	m := NewManager()
	// Should not panic
	m.Cancel(999)
}
