package techniques

import "time"

type Technique struct {
	ID          string
	Name        string
	Description string
	Duration    time.Duration
	Category    string
}

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
