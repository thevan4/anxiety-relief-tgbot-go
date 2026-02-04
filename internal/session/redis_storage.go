// Package session provides user session state management and storage.
package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStorage implements Storage interface using Redis as backend.
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage creates a new Redis storage instance and verifies connection.
func NewRedisStorage(addr string) (Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisStorage{client: client}, nil
}

const sessionTTL = 5 * time.Minute

// SetState stores user's session state in Redis with TTL.
func (r *RedisStorage) SetState(ctx context.Context, userID int64, state State) error {
	key := fmt.Sprintf("session:%d", userID)
	return r.client.Set(ctx, key, string(state), sessionTTL).Err()
}

// GetState retrieves user's session state from Redis. Returns StateUnknown if not found.
func (r *RedisStorage) GetState(ctx context.Context, userID int64) (State, error) {
	key := fmt.Sprintf("session:%d", userID)
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return StateUnknown, nil
	}
	if err != nil {
		return StateUnknown, err
	}
	return State(result), nil
}

// ClearState removes user's session state and message ID from Redis.
func (r *RedisStorage) ClearState(ctx context.Context, userID int64) error {
	stateKey := fmt.Sprintf("session:%d", userID)
	msgKey := fmt.Sprintf("message:%d", userID)
	return r.client.Del(ctx, stateKey, msgKey).Err()
}

// SetMessageID stores Telegram message ID for user session with TTL.
func (r *RedisStorage) SetMessageID(ctx context.Context, userID int64, messageID int) error {
	key := fmt.Sprintf("message:%d", userID)
	return r.client.Set(ctx, key, messageID, sessionTTL).Err()
}

// GetMessageID retrieves stored Telegram message ID for user. Returns 0 if not found.
func (r *RedisStorage) GetMessageID(ctx context.Context, userID int64) (int, error) {
	key := fmt.Sprintf("message:%d", userID)
	result, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return result, err
}

// SetLang stores user's language preference permanently (no TTL).
func (r *RedisStorage) SetLang(ctx context.Context, userID int64, lang string) error {
	key := fmt.Sprintf("lang:%d", userID)
	// Language preference is stored permanently (no TTL)
	return r.client.Set(ctx, key, lang, 0).Err()
}

// GetLang retrieves user's language preference. Returns empty string if not found.
func (r *RedisStorage) GetLang(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("lang:%d", userID)
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

// Close closes the Redis client connection.
func (r *RedisStorage) Close() error {
	return r.client.Close()
}
