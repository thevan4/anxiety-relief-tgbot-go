package sessioncore

import (
	"context"
	"sync"
)

// Manager manages user sessions with cancellation support.
// It allows cancelling running operations when user starts a new action.
// Thread-safe for concurrent access.
type Manager struct {
	mu      sync.RWMutex
	cancels map[int64]context.CancelFunc
}

// NewManager creates a new session manager instance.
func NewManager() *Manager {
	return &Manager{
		cancels: make(map[int64]context.CancelFunc),
	}
}

// Start creates a new context for user's session.
// If there's an existing session, it will be cancelled first.
// Returns a new context that should be used for the session.
func (m *Manager) Start(parentCtx context.Context, userID int64) context.Context {
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

// Cancel cancels the user's current session and removes it from tracking.
func (m *Manager) Cancel(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, exists := m.cancels[userID]; exists {
		cancel()
		delete(m.cancels, userID)
	}
}

// Has returns true if user has an active session being tracked.
func (m *Manager) Has(userID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.cancels[userID]
	return exists
}

// Count returns the number of active sessions.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.cancels)
}

// CancelAll cancels all active sessions (useful for graceful shutdown).
func (m *Manager) CancelAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for userID, cancel := range m.cancels {
		cancel()
		delete(m.cancels, userID)
	}
}
