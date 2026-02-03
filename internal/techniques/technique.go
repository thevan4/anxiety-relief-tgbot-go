package techniques

import "time"

type BreathingCycle struct {
	Name        string        // Short name: "Вдох", "Задержка", "Выдох"
	Emoji       string        // Emoji for the phase
	Duration    time.Duration // Duration in seconds
	Instruction string        // Full instruction text (legacy)
}

type GroundingStep struct {
	Number      int
	Title       string
	Description string
	Emoji       string
}
