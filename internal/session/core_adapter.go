package session

import (
	"context"
	"errors"
	"time"

	"github.com/thevan4/anxiety-relief-tgbot-go/internal/sessioncore"
)

const (
	resourceHolder = "holder"
	resourceMenu   = "menu"
)

// CoreStorageAdapter adapts session.Storage to sessioncore.Storage interface.
// This allows using the existing Redis storage with the platform-agnostic cleanup logic.
type CoreStorageAdapter struct {
	storage Storage
}

// NewCoreStorageAdapter creates adapter from existing session.Storage.
func NewCoreStorageAdapter(storage Storage) *CoreStorageAdapter {
	return &CoreStorageAdapter{storage: storage}
}

// GetState implements sessioncore.Storage.
func (a *CoreStorageAdapter) GetState(ctx context.Context, userID int64) (string, error) {
	state, err := a.storage.GetState(ctx, userID)
	return string(state), err
}

// SetState implements sessioncore.Storage.
func (a *CoreStorageAdapter) SetState(ctx context.Context, userID int64, state string) error {
	return a.storage.SetState(ctx, userID, State(state))
}

// ClearState implements sessioncore.Storage.
func (a *CoreStorageAdapter) ClearState(ctx context.Context, userID int64) error {
	return a.storage.ClearState(ctx, userID)
}

// GetResourceID implements sessioncore.Storage.
func (a *CoreStorageAdapter) GetResourceID(ctx context.Context, userID int64, kind string) (int, error) {
	switch kind {
	case resourceHolder:
		return a.storage.GetHolderMessageID(ctx, userID)
	case resourceMenu:
		return a.storage.GetMenuMessageID(ctx, userID)
	default:
		return a.storage.GetMenuMessageID(ctx, userID)
	}
}

// SetResourceID implements sessioncore.Storage.
func (a *CoreStorageAdapter) SetResourceID(ctx context.Context, userID int64, kind string, id int) error {
	switch kind {
	case resourceHolder:
		return a.storage.SetHolderMessageID(ctx, userID, id)
	case resourceMenu:
		return a.storage.SetMenuMessageID(ctx, userID, id)
	default:
		return a.storage.SetMenuMessageID(ctx, userID, id)
	}
}

// ClearResourceID implements sessioncore.Storage.
func (a *CoreStorageAdapter) ClearResourceID(ctx context.Context, userID int64, kind string) error {
	switch kind {
	case resourceHolder:
		return a.storage.ClearHolderMessageID(ctx, userID)
	case resourceMenu:
		return a.storage.ClearMenuMessageID(ctx, userID)
	default:
		return a.storage.ClearMenuMessageID(ctx, userID)
	}
}

// GetResourceCreatedAt implements sessioncore.Storage.
func (a *CoreStorageAdapter) GetResourceCreatedAt(ctx context.Context, userID int64, _ string) (time.Time, error) {
	// Currently only menu has creation timestamp
	return a.storage.GetMenuCreatedAt(ctx, userID)
}

// SetResourceCreatedAt implements sessioncore.Storage.
func (a *CoreStorageAdapter) SetResourceCreatedAt(ctx context.Context, userID int64, _ string, t time.Time) error {
	return a.storage.SetMenuCreatedAt(ctx, userID, t)
}

// AddToCleanupQueue implements sessioncore.Storage.
func (a *CoreStorageAdapter) AddToCleanupQueue(ctx context.Context, userID int64, checkAt time.Time) error {
	return a.storage.AddToCleanupQueue(ctx, userID, checkAt)
}

// GetPendingCleanup implements sessioncore.Storage.
func (a *CoreStorageAdapter) GetPendingCleanup(ctx context.Context, before time.Time, limit int64) ([]int64, error) {
	return a.storage.GetPendingCleanup(ctx, before, limit)
}

// RemoveFromCleanupQueue implements sessioncore.Storage.
func (a *CoreStorageAdapter) RemoveFromCleanupQueue(ctx context.Context, userID int64) error {
	return a.storage.RemoveFromCleanupQueue(ctx, userID)
}

// GetCleanupRetry implements sessioncore.Storage.
func (a *CoreStorageAdapter) GetCleanupRetry(ctx context.Context, userID int64) (*sessioncore.CleanupRetry, error) {
	retry, err := a.storage.GetCleanupRetry(ctx, userID)
	if errors.Is(err, ErrCleanupRetryNotFound) {
		return nil, sessioncore.ErrCleanupRetryNotFound
	}
	if err != nil || retry == nil {
		return nil, err
	}
	return &sessioncore.CleanupRetry{
		State:   retry.State,
		Attempt: retry.Attempt,
	}, nil
}

// SetCleanupRetry implements sessioncore.Storage.
func (a *CoreStorageAdapter) SetCleanupRetry(ctx context.Context, userID int64, retry sessioncore.CleanupRetry) error {
	return a.storage.SetCleanupRetry(ctx, userID, CleanupRetry{
		State:   retry.State,
		Attempt: retry.Attempt,
	})
}

// ClearCleanupRetry implements sessioncore.Storage.
func (a *CoreStorageAdapter) ClearCleanupRetry(ctx context.Context, userID int64) error {
	return a.storage.ClearCleanupRetry(ctx, userID)
}

// Ensure CoreStorageAdapter implements sessioncore.Storage.
var _ sessioncore.Storage = (*CoreStorageAdapter)(nil)
