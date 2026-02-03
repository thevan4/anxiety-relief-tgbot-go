package session

import (
	"context"
	"sync"
)

// SessionManager manages user sessions with cancellation support.
// It allows cancelling running exercises when user starts a new action.
type SessionManager struct {
	mu      sync.RWMutex
	cancels map[int64]context.CancelFunc
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		cancels: make(map[int64]context.CancelFunc),
	}
}

// StartSession creates a new context for user's session.
// If there's an existing session, it will be cancelled first.
// Returns a new context that should be used for the session.
func (m *SessionManager) StartSession(parentCtx context.Context, userID int64) context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Cancel existing session if any
	if cancel, exists := m.cancels[userID]; exists {
		cancel()
	}

	// Create new context with cancel
	ctx, cancel := context.WithCancel(parentCtx)
	m.cancels[userID] = cancel

	return ctx
}

// CancelSession cancels the user's current session.
func (m *SessionManager) CancelSession(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, exists := m.cancels[userID]; exists {
		cancel()
		delete(m.cancels, userID)
	}
}

// GetSessionContext returns a new context for the user.
// If no active session exists, creates one.
func (m *SessionManager) GetSessionContext(parentCtx context.Context, userID int64) context.Context {
	m.mu.RLock()
	_, exists := m.cancels[userID]
	m.mu.RUnlock()

	if !exists {
		return m.StartSession(parentCtx, userID)
	}

	// Return a child context of the parent
	ctx, cancel := context.WithCancel(parentCtx)
	m.mu.Lock()
	m.cancels[userID] = cancel
	m.mu.Unlock()
	return ctx
}
