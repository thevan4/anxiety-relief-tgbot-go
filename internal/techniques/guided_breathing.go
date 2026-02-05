package techniques

import "time"

// Breathing pattern timing constants.
const (
	breathSec3 = 3 * time.Second
	breathSec4 = 4 * time.Second
	breathSec6 = 6 * time.Second
	breathSec7 = 7 * time.Second
	breathSec8 = 8 * time.Second

	cycles4 = 4
	cycles5 = 5
	cycles6 = 6
)

// BreathingPattern defines a guided breathing exercise with specific timing.
type BreathingPattern struct {
	ID      string
	Emoji   string
	Inhale  time.Duration
	HoldIn  time.Duration
	Exhale  time.Duration
	HoldOut time.Duration
	Cycles  int
}

// BreathingPhase represents one phase in a guided breathing pattern.
type BreathingPhase struct {
	Phase    string
	Duration time.Duration
	Emoji    string
}

// GetBreathingPatterns returns available guided breathing patterns with different timings.
func GetBreathingPatterns() []BreathingPattern {
	return []BreathingPattern{
		{
			ID:      "box",
			Emoji:   "📦",
			Inhale:  breathSec4,
			HoldIn:  breathSec4,
			Exhale:  breathSec4,
			HoldOut: breathSec4,
			Cycles:  cycles6,
		},
		{
			ID:      "relaxing",
			Emoji:   "😴",
			Inhale:  breathSec4,
			HoldIn:  breathSec7,
			Exhale:  breathSec8,
			HoldOut: 0,
			Cycles:  cycles4,
		},
		{
			ID:      "energizing",
			Emoji:   "⚡",
			Inhale:  breathSec4,
			HoldIn:  breathSec4,
			Exhale:  breathSec6,
			HoldOut: 0,
			Cycles:  cycles6,
		},
		{
			ID:      "quick",
			Emoji:   "🚀",
			Inhale:  breathSec3,
			HoldIn:  breathSec3,
			Exhale:  breathSec3,
			HoldOut: 0,
			Cycles:  cycles5,
		},
	}
}

// GetBreathingPhases converts a breathing pattern into individual phases for execution.
func GetBreathingPhases(pattern BreathingPattern) []BreathingPhase {
	phases := []BreathingPhase{}

	if pattern.Inhale > 0 {
		phases = append(phases, BreathingPhase{
			Phase:    "inhale",
			Duration: pattern.Inhale,
			Emoji:    "💨",
		})
	}

	if pattern.HoldIn > 0 {
		phases = append(phases, BreathingPhase{
			Phase:    "hold_in",
			Duration: pattern.HoldIn,
			Emoji:    "⏸️",
		})
	}

	if pattern.Exhale > 0 {
		phases = append(phases, BreathingPhase{
			Phase:    "exhale",
			Duration: pattern.Exhale,
			Emoji:    "🌬️",
		})
	}

	if pattern.HoldOut > 0 {
		phases = append(phases, BreathingPhase{
			Phase:    "hold_out",
			Duration: pattern.HoldOut,
			Emoji:    "⏸️",
		})
	}

	return phases
}
