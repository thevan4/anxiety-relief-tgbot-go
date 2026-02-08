package techniques

import "time"

// Phase ID constants for breathing exercises.
const (
	PhaseInhale  = "inhale"
	PhaseHoldIn  = "hold_in"
	PhaseExhale  = "exhale"
	PhaseHoldOut = "hold_out"
	PhaseHold    = "hold"
)

// BreathingCycle represents one phase of breathing exercise (inhale/hold/exhale).
type BreathingCycle struct {
	Phase    string        // Phase ID: "inhale", "hold", "exhale"
	Emoji    string        // Emoji for the phase
	Duration time.Duration // Duration in seconds
}

// GroundingStep represents one step of the 5-4-3-2-1 grounding technique.
type GroundingStep struct {
	Number int
	Sense  string // "sight", "touch", "hearing", "smell", "taste"
	Count  int    // 5, 4, 3, 2, 1
	Emoji  string
}
