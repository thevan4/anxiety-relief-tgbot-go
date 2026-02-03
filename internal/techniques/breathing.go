package techniques

import "time"

const (
	BreathInhaleSec = 4
	BreathHoldSec   = 4
	BreathExhaleSec = 6
	BreathCycles    = 8
)

func GetBreathingCycles() []BreathingCycle {
	inhale := time.Duration(BreathInhaleSec) * time.Second
	hold := time.Duration(BreathHoldSec) * time.Second
	exhale := time.Duration(BreathExhaleSec) * time.Second

	cycleSteps := []BreathingCycle{
		{Name: "Вдох", Emoji: "🫁", Duration: inhale, Instruction: "Вдох через нос"},
		{Name: "Задержка", Emoji: "⏸️", Duration: hold, Instruction: "Задержите дыхание"},
		{Name: "Выдох", Emoji: "🌬️", Duration: exhale, Instruction: "Выдох через рот"},
	}

	result := make([]BreathingCycle, 0, BreathCycles*len(cycleSteps))
	for i := 0; i < BreathCycles; i++ {
		result = append(result, cycleSteps...)
	}
	return result
}
