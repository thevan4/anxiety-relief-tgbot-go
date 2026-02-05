package techniques

import (
	"testing"
	"time"
)

func TestGetBreathingPatterns(t *testing.T) {
	t.Parallel()

	patterns := GetBreathingPatterns()

	if len(patterns) == 0 {
		t.Error("expected at least one breathing pattern")
	}

	expectedPatterns := []string{
		"box",
		"relaxing",
		"energizing",
		"quick",
	}

	if len(patterns) != len(expectedPatterns) {
		t.Errorf("expected %d patterns, got %d", len(expectedPatterns), len(patterns))
	}

	for i, expected := range expectedPatterns {
		if patterns[i].ID != expected {
			t.Errorf("pattern %d: expected id %q, got %q", i, expected, patterns[i].ID)
		}
	}
}

func TestGetBreathingPatterns_ValidDurations(t *testing.T) {
	t.Parallel()

	patterns := GetBreathingPatterns()

	for _, p := range patterns {
		if p.Inhale <= 0 {
			t.Errorf("pattern %q: inhale duration should be positive", p.ID)
		}
		if p.Exhale <= 0 {
			t.Errorf("pattern %q: exhale duration should be positive", p.ID)
		}
		if p.Cycles <= 0 {
			t.Errorf("pattern %q: cycles should be positive", p.ID)
		}
	}
}

func TestGetBreathingPhases(t *testing.T) {
	t.Parallel()

	pattern := BreathingPattern{
		ID:      "test",
		Inhale:  4 * time.Second,
		HoldIn:  4 * time.Second,
		Exhale:  4 * time.Second,
		HoldOut: 4 * time.Second,
		Cycles:  1,
	}

	phases := GetBreathingPhases(pattern)

	if len(phases) != 4 {
		t.Errorf("expected 4 phases for box breathing, got %d", len(phases))
	}
}

func TestGetBreathingPhases_NoHoldOut(t *testing.T) {
	t.Parallel()

	pattern := BreathingPattern{
		ID:      "test",
		Inhale:  4 * time.Second,
		HoldIn:  7 * time.Second,
		Exhale:  8 * time.Second,
		HoldOut: 0,
		Cycles:  1,
	}

	phases := GetBreathingPhases(pattern)

	if len(phases) != 3 {
		t.Errorf("expected 3 phases for 4-7-8 breathing, got %d", len(phases))
	}
}
