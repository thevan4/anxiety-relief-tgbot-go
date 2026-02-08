// Package techniques provides anxiety relief exercise implementations.
package techniques

import "time"

const (
	// BreathInhaleSec defines inhale duration for basic breathing exercise.
	BreathInhaleSec = 4
	// BreathHoldSec defines hold duration for basic breathing exercise.
	BreathHoldSec = 4
	// BreathExhaleSec defines exhale duration for basic breathing exercise.
	BreathExhaleSec = 6
	// BreathCycles defines the number of breathing cycles in basic exercise.
	BreathCycles = 8
)

// GetBreathingCycles returns the 4-4-6 breathing exercise cycles (8 repetitions).
func GetBreathingCycles() []BreathingCycle {
	inhale := time.Duration(BreathInhaleSec) * time.Second
	hold := time.Duration(BreathHoldSec) * time.Second
	exhale := time.Duration(BreathExhaleSec) * time.Second

	cycleSteps := []BreathingCycle{
		{Phase: PhaseInhale, Emoji: "🫁", Duration: inhale},
		{Phase: PhaseHold, Emoji: "⏸️", Duration: hold},
		{Phase: PhaseExhale, Emoji: "🌬️", Duration: exhale},
	}

	result := make([]BreathingCycle, 0, BreathCycles*len(cycleSteps))
	for i := 0; i < BreathCycles; i++ {
		result = append(result, cycleSteps...)
	}
	return result
}
