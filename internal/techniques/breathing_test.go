package techniques

import (
	"testing"
	"time"
)

func TestGetBreathingCycles(t *testing.T) {
	t.Parallel()
	cycles := GetBreathingCycles()

	// 8 cycles × 3 phases = 24 total
	expectedCount := BreathCycles * 3
	if len(cycles) != expectedCount {
		t.Errorf("GetBreathingCycles() returned %d cycles, want %d", len(cycles), expectedCount)
	}

	// Check first cycle has correct phases
	phases := []struct {
		phase    string
		emoji    string
		duration time.Duration
	}{
		{"inhale", "🫁", time.Duration(BreathInhaleSec) * time.Second},
		{"hold", "⏸️", time.Duration(BreathHoldSec) * time.Second},
		{"exhale", "🌬️", time.Duration(BreathExhaleSec) * time.Second},
	}

	for i, expected := range phases {
		if cycles[i].Phase != expected.phase {
			t.Errorf("cycles[%d].Phase = %q, want %q", i, cycles[i].Phase, expected.phase)
		}
		if cycles[i].Emoji != expected.emoji {
			t.Errorf("cycles[%d].Emoji = %q, want %q", i, cycles[i].Emoji, expected.emoji)
		}
		if cycles[i].Duration != expected.duration {
			t.Errorf("cycles[%d].Duration = %v, want %v", i, cycles[i].Duration, expected.duration)
		}
	}
}

func TestBreathingConstants(t *testing.T) {
	t.Parallel()

	if BreathInhaleSec != 4 {
		t.Errorf("BreathInhaleSec = %d, want 4", BreathInhaleSec)
	}
	if BreathHoldSec != 4 {
		t.Errorf("BreathHoldSec = %d, want 4", BreathHoldSec)
	}
	if BreathExhaleSec != 6 {
		t.Errorf("BreathExhaleSec = %d, want 6", BreathExhaleSec)
	}
	if BreathCycles != 8 {
		t.Errorf("BreathCycles = %d, want 8", BreathCycles)
	}
}

func TestBreathingCyclePattern(t *testing.T) {
	t.Parallel()
	cycles := GetBreathingCycles()

	// Verify the pattern repeats correctly
	for i := 0; i < BreathCycles; i++ {
		base := i * 3
		if cycles[base].Phase != "inhale" {
			t.Errorf("cycle %d: expected inhale at position %d, got %s", i+1, base, cycles[base].Phase)
		}
		if cycles[base+1].Phase != "hold" {
			t.Errorf("cycle %d: expected hold at position %d, got %s", i+1, base+1, cycles[base+1].Phase)
		}
		if cycles[base+2].Phase != "exhale" {
			t.Errorf("cycle %d: expected exhale at position %d, got %s", i+1, base+2, cycles[base+2].Phase)
		}
	}
}
