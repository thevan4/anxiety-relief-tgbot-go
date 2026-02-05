package sessioncore

import (
	"context"
	"time"
)

// Storage defines the interface for persisting session data.
// Implementation can be Redis, memory, SQL, etc.
type Storage interface {
	// State management (short TTL, refreshed on activity).
	GetState(ctx context.Context, userID int64) (string, error)
	SetState(ctx context.Context, userID int64, state string) error
	ClearState(ctx context.Context, userID int64) error

	// Resource tracking (message IDs, etc.).
	// kind allows tracking multiple resources per user (e.g., "holder", "menu").
	GetResourceID(ctx context.Context, userID int64, kind string) (int, error)
	SetResourceID(ctx context.Context, userID int64, kind string, id int) error
	ClearResourceID(ctx context.Context, userID int64, kind string) error

	// Resource creation timestamp (for age calculation).
	GetResourceCreatedAt(ctx context.Context, userID int64, kind string) (time.Time, error)
	SetResourceCreatedAt(ctx context.Context, userID int64, kind string, t time.Time) error

	// Cleanup queue (scheduled checks).
	AddToCleanupQueue(ctx context.Context, userID int64, checkAt time.Time) error
	GetPendingCleanup(ctx context.Context, before time.Time, limit int64) ([]int64, error)
	RemoveFromCleanupQueue(ctx context.Context, userID int64) error

	// Cleanup retry state.
	GetCleanupRetry(ctx context.Context, userID int64) (*CleanupRetry, error)
	SetCleanupRetry(ctx context.Context, userID int64, retry CleanupRetry) error
	ClearCleanupRetry(ctx context.Context, userID int64) error
}

// StorageAdapter wraps platform-specific storage to implement sessioncore.Storage.
// Use this if your existing storage has different method signatures.
type StorageAdapter struct {
	GetStateFn               func(ctx context.Context, userID int64) (string, error)
	SetStateFn               func(ctx context.Context, userID int64, state string) error
	ClearStateFn             func(ctx context.Context, userID int64) error
	GetResourceIDFn          func(ctx context.Context, userID int64, kind string) (int, error)
	SetResourceIDFn          func(ctx context.Context, userID int64, kind string, id int) error
	ClearResourceIDFn        func(ctx context.Context, userID int64, kind string) error
	GetResourceCreatedAtFn   func(ctx context.Context, userID int64, kind string) (time.Time, error)
	SetResourceCreatedAtFn   func(ctx context.Context, userID int64, kind string, t time.Time) error
	AddToCleanupQueueFn      func(ctx context.Context, userID int64, checkAt time.Time) error
	GetPendingCleanupFn      func(ctx context.Context, before time.Time, limit int64) ([]int64, error)
	RemoveFromCleanupQueueFn func(ctx context.Context, userID int64) error
	GetCleanupRetryFn        func(ctx context.Context, userID int64) (*CleanupRetry, error)
	SetCleanupRetryFn        func(ctx context.Context, userID int64, retry CleanupRetry) error
	ClearCleanupRetryFn      func(ctx context.Context, userID int64) error
}

func (a *StorageAdapter) GetState(ctx context.Context, userID int64) (string, error) {
	return a.GetStateFn(ctx, userID)
}

func (a *StorageAdapter) SetState(ctx context.Context, userID int64, state string) error {
	return a.SetStateFn(ctx, userID, state)
}

func (a *StorageAdapter) ClearState(ctx context.Context, userID int64) error {
	return a.ClearStateFn(ctx, userID)
}

func (a *StorageAdapter) GetResourceID(ctx context.Context, userID int64, kind string) (int, error) {
	return a.GetResourceIDFn(ctx, userID, kind)
}

func (a *StorageAdapter) SetResourceID(ctx context.Context, userID int64, kind string, id int) error {
	return a.SetResourceIDFn(ctx, userID, kind, id)
}

func (a *StorageAdapter) ClearResourceID(ctx context.Context, userID int64, kind string) error {
	return a.ClearResourceIDFn(ctx, userID, kind)
}

func (a *StorageAdapter) GetResourceCreatedAt(ctx context.Context, userID int64, kind string) (time.Time, error) {
	return a.GetResourceCreatedAtFn(ctx, userID, kind)
}

func (a *StorageAdapter) SetResourceCreatedAt(ctx context.Context, userID int64, kind string, t time.Time) error {
	return a.SetResourceCreatedAtFn(ctx, userID, kind, t)
}

func (a *StorageAdapter) AddToCleanupQueue(ctx context.Context, userID int64, checkAt time.Time) error {
	return a.AddToCleanupQueueFn(ctx, userID, checkAt)
}

func (a *StorageAdapter) GetPendingCleanup(ctx context.Context, before time.Time, limit int64) ([]int64, error) {
	return a.GetPendingCleanupFn(ctx, before, limit)
}

func (a *StorageAdapter) RemoveFromCleanupQueue(ctx context.Context, userID int64) error {
	return a.RemoveFromCleanupQueueFn(ctx, userID)
}

func (a *StorageAdapter) GetCleanupRetry(ctx context.Context, userID int64) (*CleanupRetry, error) {
	return a.GetCleanupRetryFn(ctx, userID)
}

func (a *StorageAdapter) SetCleanupRetry(ctx context.Context, userID int64, retry CleanupRetry) error {
	return a.SetCleanupRetryFn(ctx, userID, retry)
}

func (a *StorageAdapter) ClearCleanupRetry(ctx context.Context, userID int64) error {
	return a.ClearCleanupRetryFn(ctx, userID)
}
