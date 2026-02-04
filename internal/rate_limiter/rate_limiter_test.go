package rate_limiter

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(time.Second, 10)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
	if rl.mux == nil {
		t.Error("mux is nil")
	}
	if rl.userLimits == nil {
		t.Error("userLimits is nil")
	}
	if rl.settings.timeLimit != time.Second {
		t.Errorf("timeLimit = %v, want %v", rl.settings.timeLimit, time.Second)
	}
	if rl.settings.rateLimit != 10 {
		t.Errorf("rateLimit = %d, want %d", rl.settings.rateLimit, 10)
	}
}

func TestGetUserRequestsInfo(t *testing.T) {
	rl := NewRateLimiter(time.Second, 10)

	// First call should create new info
	info1 := rl.GetUserRequestsInfo(123)
	if info1 == nil {
		t.Fatal("GetUserRequestsInfo returned nil for new user")
	}
	if info1.RateLimiter == nil {
		t.Error("RateLimiter is nil")
	}

	// Second call should return same info
	info2 := rl.GetUserRequestsInfo(123)
	if info1 != info2 {
		t.Error("GetUserRequestsInfo returned different info for same user")
	}

	// Different user should get different info
	info3 := rl.GetUserRequestsInfo(456)
	if info1 == info3 {
		t.Error("GetUserRequestsInfo returned same info for different users")
	}
}

func TestWaitAndGo(t *testing.T) {
	rl := NewRateLimiter(time.Second, 100)
	ctx := context.Background()

	// Should not block with high rate limit
	start := time.Now()
	rl.WaitAndGo(ctx, 123)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("WaitAndGo took too long: %v", elapsed)
	}
}

func TestWaitAndGoCancelledContext(t *testing.T) {
	rl := NewRateLimiter(time.Second, 1)
	ctx, cancel := context.WithCancel(context.Background())

	// Use up the rate limit
	rl.WaitAndGo(ctx, 123)

	// Cancel context before next request
	cancel()

	// Should return immediately due to cancelled context
	start := time.Now()
	rl.WaitAndGo(ctx, 123)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("WaitAndGo with cancelled context took too long: %v", elapsed)
	}
}

func TestRateLimiterInterface(_ *testing.T) {
	var _ Limiter = NewRateLimiter(time.Second, 10)
}

func TestConcurrentAccess(_ *testing.T) {
	rl := NewRateLimiter(time.Second, 100)
	ctx := context.Background()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(userID int64) {
			for j := 0; j < 10; j++ {
				rl.GetUserRequestsInfo(userID)
				rl.WaitAndGo(ctx, userID)
			}
			done <- true
		}(int64(i))
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
