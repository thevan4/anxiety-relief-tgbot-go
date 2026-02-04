package session

import (
	"context"
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("NewSessionManager returned nil")
	}
	if sm.cancels == nil {
		t.Error("cancels map is nil")
	}
}

func TestStartSession(t *testing.T) {
	sm := NewSessionManager()
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
	sm := NewSessionManager()
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
	sm := NewSessionManager()
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

func TestCancelSessionNonExistent(_ *testing.T) {
	sm := NewSessionManager()

	// Should not panic
	sm.CancelSession(999)
}

func TestGetSessionContext(t *testing.T) {
	sm := NewSessionManager()
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
	sm := NewSessionManager()
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
	sm := NewSessionManager()
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

func TestConcurrentSessionOperations(_ *testing.T) {
	sm := NewSessionManager()
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
