package session

import (
	"context"
	"errors"
	"fmt"

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

func (r *RedisStorage) SetState(ctx context.Context, userID int64, state State) error {
	key := fmt.Sprintf("session:%d", userID)
	return r.client.Set(ctx, key, string(state), 0).Err()
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
	key := fmt.Sprintf("session:%d", userID)
	return r.client.Del(ctx, key).Err()
}
