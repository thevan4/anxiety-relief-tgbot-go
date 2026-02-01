package techniques

import "time"

const (
	breathInhaleSec = 4
	breathHoldSec   = 4
	breathExhaleSec = 6
	breathCycles    = 8
)

func GetBreathingCycles() []BreathingCycle {
	inhale := time.Duration(breathInhaleSec) * time.Second
	hold := time.Duration(breathHoldSec) * time.Second
	exhale := time.Duration(breathExhaleSec) * time.Second

	cycleSteps := []BreathingCycle{
		{Instruction: "💨 Вдох через нос (4 сек)", Duration: inhale},
		{Instruction: "⏸️ Задержите дыхание (4 сек)", Duration: hold},
		{Instruction: "🌬️ Выдох через рот (6 сек)", Duration: exhale},
	}

	result := make([]BreathingCycle, 0, breathCycles*len(cycleSteps))
	for i := 0; i < breathCycles; i++ {
		result = append(result, cycleSteps...)
	}
	return result
}
