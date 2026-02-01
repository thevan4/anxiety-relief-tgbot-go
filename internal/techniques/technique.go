package techniques

import "time"

type BreathingCycle struct {
	Instruction string
	Duration    time.Duration
}

type GroundingStep struct {
	Number      int
	Title       string
	Description string
	Emoji       string
}
