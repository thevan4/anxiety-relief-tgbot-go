package session

import (
	"context"
)

type Storage interface {
	SetState(ctx context.Context, userID int64, state State) error
	GetState(ctx context.Context, userID int64) (State, error)
	ClearState(ctx context.Context, userID int64) error
	Close() error
}
