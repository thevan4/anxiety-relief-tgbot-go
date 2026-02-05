// Package session provides user session state management and storage.
package session

import (
	"context"
	"encoding/json"
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

const (
	// SessionMaxTTL — максимальное время жизни сессии (ограничение Telegram API на edit/delete).
	SessionMaxTTL = 48 * time.Hour

	// CleanupDelay — время до первой проверки после создания меню (47 часов).
	CleanupDelay = 47 * time.Hour

	// stateTTL — TTL для state ключа (короткий, обновляется при активности).
	stateTTL = 5 * time.Minute

	// messageIDTTL — TTL для message ID ключей (равен SessionMaxTTL).
	messageIDTTL = SessionMaxTTL

	// cleanupRetryTTL — TTL для retry ключей.
	cleanupRetryTTL = 1 * time.Hour

	// cleanupQueueKey — имя ZSET для очереди очистки.
	cleanupQueueKey = "cleanup_queue"
)

// SetState stores user's session state in Redis with TTL.
func (r *RedisStorage) SetState(ctx context.Context, userID int64, state State) error {
	key := fmt.Sprintf("session:%d", userID)
	return r.client.Set(ctx, key, string(state), stateTTL).Err()
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

// ClearState removes user's session state from Redis.
// Note: messageID is NOT cleared - it's needed to validate callbacks and delete old messages.
func (r *RedisStorage) ClearState(ctx context.Context, userID int64) error {
	stateKey := fmt.Sprintf("session:%d", userID)
	return r.client.Del(ctx, stateKey).Err()
}

// SetMessageID stores Telegram message ID for user session.
//
// Deprecated: Use SetMenuMessageID instead.
func (r *RedisStorage) SetMessageID(ctx context.Context, userID int64, messageID int) error {
	return r.SetMenuMessageID(ctx, userID, messageID)
}

// GetMessageID retrieves stored Telegram message ID for user. Returns 0 if not found.
//
// Deprecated: Use GetMenuMessageID instead.
func (r *RedisStorage) GetMessageID(ctx context.Context, userID int64) (int, error) {
	return r.GetMenuMessageID(ctx, userID)
}

// SetHolderMessageID stores holder (welcome) message ID.
func (r *RedisStorage) SetHolderMessageID(ctx context.Context, userID int64, messageID int) error {
	key := fmt.Sprintf("holder:%d", userID)
	return r.client.Set(ctx, key, messageID, messageIDTTL).Err()
}

// GetHolderMessageID retrieves holder message ID. Returns 0 if not found.
func (r *RedisStorage) GetHolderMessageID(ctx context.Context, userID int64) (int, error) {
	key := fmt.Sprintf("holder:%d", userID)
	result, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return result, err
}

// ClearHolderMessageID removes holder message ID.
func (r *RedisStorage) ClearHolderMessageID(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("holder:%d", userID)
	return r.client.Del(ctx, key).Err()
}

// SetMenuMessageID stores menu/exercise message ID.
func (r *RedisStorage) SetMenuMessageID(ctx context.Context, userID int64, messageID int) error {
	key := fmt.Sprintf("menu:%d", userID)
	return r.client.Set(ctx, key, messageID, messageIDTTL).Err()
}

// GetMenuMessageID retrieves menu message ID. Returns 0 if not found.
func (r *RedisStorage) GetMenuMessageID(ctx context.Context, userID int64) (int, error) {
	key := fmt.Sprintf("menu:%d", userID)
	result, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return result, err
}

// ClearMenuMessageID removes menu message ID.
func (r *RedisStorage) ClearMenuMessageID(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("menu:%d", userID)
	return r.client.Del(ctx, key).Err()
}

// SetMenuCreatedAt stores menu creation timestamp.
func (r *RedisStorage) SetMenuCreatedAt(ctx context.Context, userID int64, t time.Time) error {
	key := fmt.Sprintf("menu_created:%d", userID)
	return r.client.Set(ctx, key, t.Unix(), messageIDTTL).Err()
}

// GetMenuCreatedAt retrieves menu creation timestamp. Returns zero time if not found.
func (r *RedisStorage) GetMenuCreatedAt(ctx context.Context, userID int64) (time.Time, error) {
	key := fmt.Sprintf("menu_created:%d", userID)
	result, err := r.client.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(result, 0), nil
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

// AddToCleanupQueue adds user to cleanup queue with scheduled check time.
func (r *RedisStorage) AddToCleanupQueue(ctx context.Context, userID int64, checkAt time.Time) error {
	return r.client.ZAdd(ctx, cleanupQueueKey, redis.Z{
		Score:  float64(checkAt.Unix()),
		Member: userID,
	}).Err()
}

// GetPendingCleanup returns user IDs that need cleanup check (checkAt <= now).
func (r *RedisStorage) GetPendingCleanup(ctx context.Context, now time.Time, limit int64) ([]int64, error) {
	results, err := r.client.ZRangeByScore(ctx, cleanupQueueKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   fmt.Sprintf("%d", now.Unix()),
		Count: limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, len(results))
	for _, s := range results {
		var userID int64
		if _, err := fmt.Sscanf(s, "%d", &userID); err == nil {
			userIDs = append(userIDs, userID)
		}
	}
	return userIDs, nil
}

// RemoveFromCleanupQueue removes user from cleanup queue.
func (r *RedisStorage) RemoveFromCleanupQueue(ctx context.Context, userID int64) error {
	return r.client.ZRem(ctx, cleanupQueueKey, userID).Err()
}

// SetCleanupRetry stores retry information for user.
func (r *RedisStorage) SetCleanupRetry(ctx context.Context, userID int64, retry CleanupRetry) error {
	key := fmt.Sprintf("cleanup_retry:%d", userID)
	data, err := json.Marshal(retry)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, cleanupRetryTTL).Err()
}

// GetCleanupRetry retrieves retry information for user. Returns nil if not found.
func (r *RedisStorage) GetCleanupRetry(ctx context.Context, userID int64) (*CleanupRetry, error) {
	key := fmt.Sprintf("cleanup_retry:%d", userID)
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCleanupRetryNotFound
	}
	if err != nil {
		return nil, err
	}

	var retry CleanupRetry
	if err := json.Unmarshal([]byte(result), &retry); err != nil {
		return nil, err
	}
	return &retry, nil
}

// ClearCleanupRetry removes retry information for user.
func (r *RedisStorage) ClearCleanupRetry(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("cleanup_retry:%d", userID)
	return r.client.Del(ctx, key).Err()
}

// ClearSession removes all session data (state, messages).
// Language preference is NOT cleared.
func (r *RedisStorage) ClearSession(ctx context.Context, userID int64) error {
	keys := []string{
		fmt.Sprintf("session:%d", userID),
		fmt.Sprintf("holder:%d", userID),
		fmt.Sprintf("menu:%d", userID),
		fmt.Sprintf("menu_created:%d", userID),
		fmt.Sprintf("cleanup_retry:%d", userID),
	}
	return r.client.Del(ctx, keys...).Err()
}

// Close closes the Redis client connection.
func (r *RedisStorage) Close() error {
	return r.client.Close()
}
