package session

import (
	"context"
)

type Storage interface {
	SetState(ctx context.Context, userID int64, state State) error
	GetState(ctx context.Context, userID int64) (State, error)
	ClearState(ctx context.Context, userID int64) error

	// SetMessageID stores bot's message ID for editing later
	SetMessageID(ctx context.Context, userID int64, messageID int) error
	GetMessageID(ctx context.Context, userID int64) (int, error)

	// SetLang stores user's language preference
	SetLang(ctx context.Context, userID int64, lang string) error
	// GetLang returns user's language preference (empty string if not set)
	GetLang(ctx context.Context, userID int64) (string, error)

	Close() error
}
