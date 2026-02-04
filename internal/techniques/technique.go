package techniques

import "time"

// BreathingCycle represents one phase of breathing exercise (inhale/hold/exhale).
type BreathingCycle struct {
	Name        string        // Short name: "Вдох", "Задержка", "Выдох"
	Emoji       string        // Emoji for the phase
	Duration    time.Duration // Duration in seconds
	Instruction string        // Full instruction text (legacy)
}

// GroundingStep represents one step of the 5-4-3-2-1 grounding technique.
type GroundingStep struct {
	Number      int
	Title       string
	Description string
	Emoji       string
}
