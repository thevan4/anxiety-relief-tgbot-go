package session

import (
	"testing"
	"time"
)

func TestCleanupProcessor_EvaluateCleanup(t *testing.T) {
	t.Parallel()
	config := DefaultConfig()
	processor := NewCleanupProcessor(config)
	now := time.Now()

	tests := []struct {
		name         string
		createdAt    time.Time
		currentState string
		retry        *CleanupRetry
		wantDecision CleanupDecision
	}{
		{
			name:         "resource too old",
			createdAt:    now.Add(-49 * time.Hour),
			currentState: "some_state",
			retry:        nil,
			wantDecision: DecisionDelete,
		},
		{
			name:         "no creation timestamp",
			createdAt:    time.Time{},
			currentState: "some_state",
			retry:        nil,
			wantDecision: DecisionDelete,
		},
		{
			name:         "first check - user inactive",
			createdAt:    now.Add(-47 * time.Hour),
			currentState: "",
			retry:        nil,
			wantDecision: DecisionDelete,
		},
		{
			name:         "first check - user active",
			createdAt:    now.Add(-47 * time.Hour),
			currentState: "breathing_running",
			retry:        nil,
			wantDecision: DecisionRetry,
		},
		{
			name:         "retry - user finished (became inactive)",
			createdAt:    now.Add(-47 * time.Hour),
			currentState: "",
			retry:        &CleanupRetry{State: "breathing_running", Attempt: 1},
			wantDecision: DecisionKeep,
		},
		{
			name:         "retry - user stuck in same state",
			createdAt:    now.Add(-47 * time.Hour),
			currentState: "breathing_running",
			retry:        &CleanupRetry{State: "breathing_running", Attempt: 1},
			wantDecision: DecisionDelete,
		},
		{
			name:         "retry - user switched states",
			createdAt:    now.Add(-47 * time.Hour),
			currentState: "grounding_step_1",
			retry:        &CleanupRetry{State: "breathing_running", Attempt: 1},
			wantDecision: DecisionRetry,
		},
		{
			name:         "retry - max retries exceeded",
			createdAt:    now.Add(-47 * time.Hour),
			currentState: "grounding_step_1",
			retry:        &CleanupRetry{State: "breathing_running", Attempt: 6},
			wantDecision: DecisionDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := processor.EvaluateCleanup(tt.createdAt, tt.currentState, tt.retry, now)
			if result.Decision != tt.wantDecision {
				t.Errorf("EvaluateCleanup() decision = %v, want %v (reason: %s)",
					result.Decision, tt.wantDecision, result.Reason)
			}
		})
	}
}

func TestCleanupProcessor_ScheduleInitialCleanup(t *testing.T) {
	t.Parallel()
	config := DefaultConfig()
	processor := NewCleanupProcessor(config)

	createdAt := time.Now()
	expected := createdAt.Add(config.CleanupDelay)

	got := processor.ScheduleInitialCleanup(createdAt)
	if !got.Equal(expected) {
		t.Errorf("ScheduleInitialCleanup() = %v, want %v", got, expected)
	}
}

func TestCleanupProcessor_Config(t *testing.T) {
	t.Parallel()
	config := Config{
		StateTTL:      10 * time.Minute,
		ResourceTTL:   24 * time.Hour,
		CleanupDelay:  23 * time.Hour,
		RetryInterval: 10 * time.Minute,
		MaxRetries:    3,
	}
	processor := NewCleanupProcessor(config)

	got := processor.Config()
	if got != config {
		t.Errorf("Config() = %v, want %v", got, config)
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	config := DefaultConfig()

	if config.StateTTL != 5*time.Minute {
		t.Errorf("StateTTL = %v, want 5m", config.StateTTL)
	}
	if config.ResourceTTL != 48*time.Hour {
		t.Errorf("ResourceTTL = %v, want 48h", config.ResourceTTL)
	}
	if config.CleanupDelay != 47*time.Hour {
		t.Errorf("CleanupDelay = %v, want 47h", config.CleanupDelay)
	}
	if config.MaxRetries != 6 {
		t.Errorf("MaxRetries = %v, want 6", config.MaxRetries)
	}
}
