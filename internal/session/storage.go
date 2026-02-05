package session

import (
	"context"
	"errors"
	"time"
)

// ErrCleanupRetryNotFound indicates missing cleanup retry entry.
//
//nolint:gochecknoglobals // Sentinel error for storage.
var ErrCleanupRetryNotFound = errors.New("cleanup retry not found")

// CleanupRetry stores retry information for cleanup queue.
type CleanupRetry struct {
	State   string `json:"state"`
	Attempt int    `json:"attempt"`
}

// Storage defines the interface for persisting user session data.
type Storage interface {
	// State management.
	SetState(ctx context.Context, userID int64, state State) error
	GetState(ctx context.Context, userID int64) (State, error)
	ClearState(ctx context.Context, userID int64) error

	// Message IDs (legacy - kept for backward compatibility).
	//
	// Deprecated: Use SetHolderMessageID/SetMenuMessageID instead.
	SetMessageID(ctx context.Context, userID int64, messageID int) error
	GetMessageID(ctx context.Context, userID int64) (int, error)

	// Holder message (приветствие после /start).
	SetHolderMessageID(ctx context.Context, userID int64, messageID int) error
	GetHolderMessageID(ctx context.Context, userID int64) (int, error)
	ClearHolderMessageID(ctx context.Context, userID int64) error

	// Menu message (меню/упражнение после "Начать").
	SetMenuMessageID(ctx context.Context, userID int64, messageID int) error
	GetMenuMessageID(ctx context.Context, userID int64) (int, error)
	ClearMenuMessageID(ctx context.Context, userID int64) error

	// Menu creation timestamp (для определения возраста меню).
	SetMenuCreatedAt(ctx context.Context, userID int64, t time.Time) error
	GetMenuCreatedAt(ctx context.Context, userID int64) (time.Time, error)

	// Language preference.
	SetLang(ctx context.Context, userID int64, lang string) error
	GetLang(ctx context.Context, userID int64) (string, error)

	// Cleanup queue (ZSET: userID → timestamp проверки).
	AddToCleanupQueue(ctx context.Context, userID int64, checkAt time.Time) error
	GetPendingCleanup(ctx context.Context, now time.Time, limit int64) ([]int64, error)
	RemoveFromCleanupQueue(ctx context.Context, userID int64) error

	// Cleanup retry (для повторных проверок активных пользователей).
	SetCleanupRetry(ctx context.Context, userID int64, retry CleanupRetry) error
	GetCleanupRetry(ctx context.Context, userID int64) (*CleanupRetry, error)
	ClearCleanupRetry(ctx context.Context, userID int64) error

	// ClearSession removes all session data (state, messages).
	ClearSession(ctx context.Context, userID int64) error

	Close() error
}
