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

type BreathingPattern struct {
	Name        string
	Description string
	Emoji       string
	Inhale      time.Duration
	HoldIn      time.Duration
	Exhale      time.Duration
	HoldOut     time.Duration
	Cycles      int
}

type BreathingPhase struct {
	Name     string
	Duration time.Duration
	Emoji    string
}

func GetBreathingPatterns() []BreathingPattern {
	return []BreathingPattern{
		{
			Name:        "Коробочное дыхание",
			Description: "4-4-4-4 — для фокуса и концентрации",
			Emoji:       "📦",
			Inhale:      breathSec4,
			HoldIn:      breathSec4,
			Exhale:      breathSec4,
			HoldOut:     breathSec4,
			Cycles:      cycles6,
		},
		{
			Name:        "Расслабляющее 4-7-8",
			Description: "Для глубокого расслабления и сна",
			Emoji:       "😴",
			Inhale:      breathSec4,
			HoldIn:      breathSec7,
			Exhale:      breathSec8,
			HoldOut:     0,
			Cycles:      cycles4,
		},
		{
			Name:        "Энергичное 4-4-6",
			Description: "Для бодрости и ясности ума",
			Emoji:       "⚡",
			Inhale:      breathSec4,
			HoldIn:      breathSec4,
			Exhale:      breathSec6,
			HoldOut:     0,
			Cycles:      cycles6,
		},
		{
			Name:        "Быстрый сброс 3-3-3",
			Description: "Быстрое снятие напряжения",
			Emoji:       "🚀",
			Inhale:      breathSec3,
			HoldIn:      breathSec3,
			Exhale:      breathSec3,
			HoldOut:     0,
			Cycles:      cycles5,
		},
	}
}

func GetBreathingPhases(pattern BreathingPattern) []BreathingPhase {
	phases := []BreathingPhase{}

	if pattern.Inhale > 0 {
		phases = append(phases, BreathingPhase{
			Name:     "Вдох через нос",
			Duration: pattern.Inhale,
			Emoji:    "💨",
		})
	}

	if pattern.HoldIn > 0 {
		phases = append(phases, BreathingPhase{
			Name:     "Задержка",
			Duration: pattern.HoldIn,
			Emoji:    "⏸️",
		})
	}

	if pattern.Exhale > 0 {
		phases = append(phases, BreathingPhase{
			Name:     "Выдох через рот",
			Duration: pattern.Exhale,
			Emoji:    "🌬️",
		})
	}

	if pattern.HoldOut > 0 {
		phases = append(phases, BreathingPhase{
			Name:     "Пауза",
			Duration: pattern.HoldOut,
			Emoji:    "⏸️",
		})
	}

	return phases
}
