package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

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

func (r *RedisStorage) SetState(ctx context.Context, userID int64, state State) error {
	key := fmt.Sprintf("session:%d", userID)
	return r.client.Set(ctx, key, string(state), sessionTTL).Err()
}

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

func (r *RedisStorage) ClearState(ctx context.Context, userID int64) error {
	stateKey := fmt.Sprintf("session:%d", userID)
	msgKey := fmt.Sprintf("message:%d", userID)
	return r.client.Del(ctx, stateKey, msgKey).Err()
}

func (r *RedisStorage) SetMessageID(ctx context.Context, userID int64, messageID int) error {
	key := fmt.Sprintf("message:%d", userID)
	return r.client.Set(ctx, key, messageID, sessionTTL).Err()
}

func (r *RedisStorage) GetMessageID(ctx context.Context, userID int64) (int, error) {
	key := fmt.Sprintf("message:%d", userID)
	result, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return result, err
}

func (r *RedisStorage) SetLang(ctx context.Context, userID int64, lang string) error {
	key := fmt.Sprintf("lang:%d", userID)
	// Language preference is stored permanently (no TTL)
	return r.client.Set(ctx, key, lang, 0).Err()
}

func (r *RedisStorage) GetLang(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("lang:%d", userID)
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
