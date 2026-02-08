package session

import (
	"testing"
	"time"
)

func TestCleanupProcessor_EvaluateCleanup_Initial(t *testing.T) {
	t.Parallel()
	config := DefaultConfig()
	processor := NewCleanupProcessor(config)
	now := time.Now()

	tests := []struct {
		name         string
		createdAt    time.Time
		currentState string
		wantDecision CleanupDecision
	}{
		{"resource too old", now.Add(-49 * time.Hour), "some_state", DecisionDelete},
		{"no creation timestamp", time.Time{}, "some_state", DecisionDelete},
		{"first check - user inactive", now.Add(-47 * time.Hour), "", DecisionDelete},
		{"first check - user active", now.Add(-47 * time.Hour), "breathing_running", DecisionRetry},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := processor.EvaluateCleanup(tt.createdAt, tt.currentState, nil, now)
			if result.Decision != tt.wantDecision {
				t.Errorf("EvaluateCleanup() decision = %v, want %v (reason: %s)",
					result.Decision, tt.wantDecision, result.Reason)
			}
		})
	}
}

func TestCleanupProcessor_EvaluateCleanup_Retry(t *testing.T) {
	t.Parallel()
	config := DefaultConfig()
	processor := NewCleanupProcessor(config)
	now := time.Now()
	created := now.Add(-47 * time.Hour)

	tests := []struct {
		name         string
		currentState string
		retry        *CleanupRetry
		wantDecision CleanupDecision
	}{
		{"user finished", "", &CleanupRetry{State: "breathing_running", Attempt: 1}, DecisionKeep},
		{"user stuck", "breathing_running", &CleanupRetry{State: "breathing_running", Attempt: 1}, DecisionDelete},
		{"user switched", "grounding_step_1", &CleanupRetry{State: "breathing_running", Attempt: 1}, DecisionRetry},
		{"max retries", "grounding_step_1", &CleanupRetry{State: "breathing_running", Attempt: 6}, DecisionDelete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := processor.EvaluateCleanup(created, tt.currentState, tt.retry, now)
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
