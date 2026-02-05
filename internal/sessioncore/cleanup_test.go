package sessioncore

import (
	"testing"
	"time"
)

func TestCleanupProcessor_EvaluateCleanup(t *testing.T) {
	cfg := DefaultConfig()
	p := NewCleanupProcessor(cfg)
	now := time.Now()

	tests := []struct {
		name              string
		resourceCreatedAt time.Time
		currentState      string
		retry             *CleanupRetry
		wantDecision      CleanupDecision
		wantReason        string
	}{
		// Edge cases for resource creation time
		{
			name:              "no creation time - delete",
			resourceCreatedAt: time.Time{},
			currentState:      "",
			retry:             nil,
			wantDecision:      DecisionDelete,
			wantReason:        "no resource creation timestamp",
		},
		{
			name:              "no creation time even if user active - delete",
			resourceCreatedAt: time.Time{},
			currentState:      "active_state",
			retry:             nil,
			wantDecision:      DecisionDelete,
			wantReason:        "no resource creation timestamp",
		},
		{
			name:              "resource too old - delete",
			resourceCreatedAt: now.Add(-49 * time.Hour),
			currentState:      "active",
			retry:             nil,
			wantDecision:      DecisionDelete,
			wantReason:        "resource exceeded TTL",
		},
		{
			name:              "resource exactly at TTL - delete",
			resourceCreatedAt: now.Add(-48*time.Hour - 1*time.Second),
			currentState:      "active",
			retry:             nil,
			wantDecision:      DecisionDelete,
			wantReason:        "resource exceeded TTL",
		},
		{
			name:              "resource just under TTL - proceed with check",
			resourceCreatedAt: now.Add(-47*time.Hour - 59*time.Minute),
			currentState:      "",
			retry:             nil,
			wantDecision:      DecisionDelete,
			wantReason:        "user inactive on first check",
		},

		// First check scenarios (retry == nil)
		{
			name:              "first check, user inactive (empty state) - delete",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "",
			retry:             nil,
			wantDecision:      DecisionDelete,
			wantReason:        "user inactive on first check",
		},
		{
			name:              "first check, user active - retry",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "exercise_running",
			retry:             nil,
			wantDecision:      DecisionRetry,
			wantReason:        "user active, scheduling retry",
		},

		// Retry check: user finished (state became empty)
		{
			name:              "retry, user finished activity - keep",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "",
			retry:             &CleanupRetry{State: "exercise_running", Attempt: 1},
			wantDecision:      DecisionKeep,
			wantReason:        "user finished activity, keeping resource",
		},
		{
			name:              "retry attempt 3, user finished - keep",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "",
			retry:             &CleanupRetry{State: "some_state", Attempt: 3},
			wantDecision:      DecisionKeep,
			wantReason:        "user finished activity, keeping resource",
		},

		// Retry check: user stuck in same state
		{
			name:              "retry, user stuck in same state - delete",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "exercise_running",
			retry:             &CleanupRetry{State: "exercise_running", Attempt: 1},
			wantDecision:      DecisionDelete,
			wantReason:        "user stuck in same state",
		},
		{
			name:              "retry attempt 5, user stuck - delete",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "same_state",
			retry:             &CleanupRetry{State: "same_state", Attempt: 5},
			wantDecision:      DecisionDelete,
			wantReason:        "user stuck in same state",
		},

		// Retry check: user switched states (active)
		{
			name:              "retry, user switched state, attempt 1 - retry again",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "another_exercise",
			retry:             &CleanupRetry{State: "exercise_running", Attempt: 1},
			wantDecision:      DecisionRetry,
			wantReason:        "user switched states, scheduling retry",
		},
		{
			name:              "retry, user switched state, attempt 5 - retry again",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "new_exercise",
			retry:             &CleanupRetry{State: "old_exercise", Attempt: 5},
			wantDecision:      DecisionRetry,
			wantReason:        "user switched states, scheduling retry",
		},

		// Retry check: max retries exceeded
		{
			name:              "retry, max attempts reached (6) - delete",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "new_exercise",
			retry:             &CleanupRetry{State: "old_exercise", Attempt: 6},
			wantDecision:      DecisionDelete,
			wantReason:        "max retries exceeded",
		},
		{
			name:              "retry, over max attempts (7) - delete",
			resourceCreatedAt: now.Add(-47 * time.Hour),
			currentState:      "new_state",
			retry:             &CleanupRetry{State: "other_state", Attempt: 7},
			wantDecision:      DecisionDelete,
			wantReason:        "max retries exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.EvaluateCleanup(tt.resourceCreatedAt, tt.currentState, tt.retry, now)
			if result.Decision != tt.wantDecision {
				t.Errorf("got decision %v, want %v (reason: %s)", result.Decision, tt.wantDecision, result.Reason)
			}
			if tt.wantReason != "" && result.Reason != tt.wantReason {
				t.Errorf("got reason %q, want %q", result.Reason, tt.wantReason)
			}
		})
	}
}

func TestCleanupProcessor_RetryIncrement(t *testing.T) {
	cfg := DefaultConfig()
	p := NewCleanupProcessor(cfg)
	now := time.Now()
	resourceCreated := now.Add(-47 * time.Hour)

	// Simulate full retry flow with state changes

	// Step 1: First check when active
	result := p.EvaluateCleanup(resourceCreated, "state_a", nil, now)
	if result.Decision != DecisionRetry {
		t.Fatalf("step 1: expected retry, got %v", result.Decision)
	}
	if result.Retry == nil || result.Retry.Attempt != 1 {
		t.Fatalf("step 1: expected attempt 1, got %v", result.Retry)
	}
	if result.Retry.State != "state_a" {
		t.Errorf("step 1: expected state 'state_a', got %q", result.Retry.State)
	}

	// Step 2: User switches to state_b
	result = p.EvaluateCleanup(resourceCreated, "state_b", result.Retry, now)
	if result.Decision != DecisionRetry {
		t.Fatalf("step 2: expected retry, got %v", result.Decision)
	}
	if result.Retry.Attempt != 2 {
		t.Errorf("step 2: expected attempt 2, got %d", result.Retry.Attempt)
	}
	if result.Retry.State != "state_b" {
		t.Errorf("step 2: expected state 'state_b', got %q", result.Retry.State)
	}

	// Step 3: User switches to state_c
	result = p.EvaluateCleanup(resourceCreated, "state_c", result.Retry, now)
	if result.Decision != DecisionRetry {
		t.Fatalf("step 3: expected retry, got %v", result.Decision)
	}
	if result.Retry.Attempt != 3 {
		t.Errorf("step 3: expected attempt 3, got %d", result.Retry.Attempt)
	}

	// Step 4-6: Continue switching until max
	for i := 4; i <= 6; i++ {
		result = p.EvaluateCleanup(resourceCreated, "state_"+string(rune('a'+i)), result.Retry, now)
		if result.Decision != DecisionRetry {
			t.Fatalf("step %d: expected retry, got %v", i, result.Decision)
		}
	}

	// Step 7: Max retries, should delete
	result = p.EvaluateCleanup(resourceCreated, "state_final", result.Retry, now)
	if result.Decision != DecisionDelete {
		t.Errorf("step 7: expected delete after max retries, got %v", result.Decision)
	}
}

func TestCleanupProcessor_RetryNextCheckTime(t *testing.T) {
	cfg := Config{
		StateTTL:      5 * time.Minute,
		ResourceTTL:   48 * time.Hour,
		CleanupDelay:  47 * time.Hour,
		RetryInterval: 7 * time.Minute, // Custom interval
		MaxRetries:    6,
	}
	p := NewCleanupProcessor(cfg)
	now := time.Now()

	result := p.EvaluateCleanup(
		now.Add(-47*time.Hour),
		"active_state",
		nil,
		now,
	)

	if result.Decision != DecisionRetry {
		t.Fatalf("expected retry, got %v", result.Decision)
	}

	expectedNext := now.Add(7 * time.Minute)
	if !result.NextCheck.Equal(expectedNext) {
		t.Errorf("expected next check at %v, got %v", expectedNext, result.NextCheck)
	}
}

func TestCleanupProcessor_ScheduleInitialCleanup(t *testing.T) {
	cfg := DefaultConfig()
	p := NewCleanupProcessor(cfg)

	created := time.Now()
	scheduled := p.ScheduleInitialCleanup(created)

	expected := created.Add(cfg.CleanupDelay)
	if !scheduled.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, scheduled)
	}
}

func TestCleanupProcessor_CustomConfig(t *testing.T) {
	cfg := Config{
		StateTTL:      1 * time.Minute,
		ResourceTTL:   24 * time.Hour, // Shorter TTL
		CleanupDelay:  23 * time.Hour,
		RetryInterval: 2 * time.Minute,
		MaxRetries:    3, // Fewer retries
	}
	p := NewCleanupProcessor(cfg)
	now := time.Now()

	// Test with custom TTL
	result := p.EvaluateCleanup(
		now.Add(-25*time.Hour), // Over 24h custom TTL
		"active",
		nil,
		now,
	)
	if result.Decision != DecisionDelete {
		t.Errorf("expected delete for resource over custom TTL, got %v", result.Decision)
	}

	// Test with custom max retries
	result = p.EvaluateCleanup(
		now.Add(-23*time.Hour),
		"new_state",
		&CleanupRetry{State: "old_state", Attempt: 3}, // At max for this config
		now,
	)
	if result.Decision != DecisionDelete {
		t.Errorf("expected delete at custom max retries, got %v", result.Decision)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.StateTTL != 5*time.Minute {
		t.Errorf("StateTTL: got %v, want 5m", cfg.StateTTL)
	}
	if cfg.ResourceTTL != 48*time.Hour {
		t.Errorf("ResourceTTL: got %v, want 48h", cfg.ResourceTTL)
	}
	if cfg.CleanupDelay != 47*time.Hour {
		t.Errorf("CleanupDelay: got %v, want 47h", cfg.CleanupDelay)
	}
	if cfg.RetryInterval != 5*time.Minute {
		t.Errorf("RetryInterval: got %v, want 5m", cfg.RetryInterval)
	}
	if cfg.MaxRetries != 6 {
		t.Errorf("MaxRetries: got %d, want 6", cfg.MaxRetries)
	}
}
