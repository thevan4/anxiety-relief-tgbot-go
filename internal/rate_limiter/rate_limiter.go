package rate_limiter

import (
	"context"
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limiter ...
type Limiter interface {
	GetUserRequestsInfo(userID int64) *UserRequestsInfo
	WaitAndGo(ctx context.Context, userID int64)
}

// RateLimiter ...
type RateLimiter struct {
	mux        *sync.Mutex
	settings   setting
	userLimits map[int64]*UserRequestsInfo
}

type setting struct {
	timeLimit time.Duration
	rateLimit int
}

// UserRequestsInfo ...
type UserRequestsInfo struct {
	mux         *sync.Mutex
	RateLimiter *rate.Limiter
}

// NewRateLimiter ...
func NewRateLimiter(
	timeLimit time.Duration,
	rateLimit int,
) *RateLimiter {
	return &RateLimiter{
		mux: new(sync.Mutex),
		settings: setting{
			timeLimit: timeLimit,
			rateLimit: rateLimit,
		},
		userLimits: make(map[int64]*UserRequestsInfo),
	}
}

// GetUserRequestsInfo ...
func (rl *RateLimiter) GetUserRequestsInfo(userID int64) *UserRequestsInfo {
	rl.mux.Lock()
	defer rl.mux.Unlock()
	userRequestsInfo, exists := rl.userLimits[userID]
	if !exists {
		interval := rl.settings.timeLimit / time.Duration(rl.settings.rateLimit)
		userRequestsInfo = &UserRequestsInfo{
			mux:         new(sync.Mutex),
			RateLimiter: rate.NewLimiter(rate.Every(interval), rl.settings.rateLimit),
		}
		rl.userLimits[userID] = userRequestsInfo

		return userRequestsInfo
	}

	return userRequestsInfo
}

// WaitAndGo ожидает если надо по рейтлимиту, либо идёт дальше.
func (rl *RateLimiter) WaitAndGo(ctx context.Context, userID int64) {
	userRequestsInfo := rl.GetUserRequestsInfo(userID)
	// Blocks until the next event is resolved.
	err := userRequestsInfo.RateLimiter.Wait(ctx)
	if err != nil {
		log.Printf("rate limiter: context done for user %d: %v", userID, err)
		return
	}
	// go!
}
