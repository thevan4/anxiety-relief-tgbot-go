package messages

import (
	"strings"
	"testing"
)

func TestBreathingCycleText(t *testing.T) {
	result := BreathingCycleText(3, 8, "Вдох через нос")

	if !strings.Contains(result, "Цикл 3 из 8") {
		t.Errorf("expected cycle info, got %q", result)
	}
	if !strings.Contains(result, "Вдох через нос") {
		t.Errorf("expected instruction, got %q", result)
	}
}

func TestBreathingPhaseText(t *testing.T) {
	tests := []struct {
		name        string
		cycleNum    int
		totalCycles int
		phaseName   string
		emoji       string
		elapsed     int
		total       int
		wantCycle   string
		wantPhase   string
	}{
		{
			name:        "inhale start",
			cycleNum:    1,
			totalCycles: 8,
			phaseName:   "Вдох",
			emoji:       "🫁",
			elapsed:     0,
			total:       4,
			wantCycle:   "Цикл 1/8",
			wantPhase:   "Вдох 🫁",
		},
		{
			name:        "hold middle",
			cycleNum:    4,
			totalCycles: 8,
			phaseName:   "Задержка",
			emoji:       "⏸️",
			elapsed:     2,
			total:       4,
			wantCycle:   "Цикл 4/8",
			wantPhase:   "Задержка ⏸️",
		},
		{
			name:        "exhale end",
			cycleNum:    8,
			totalCycles: 8,
			phaseName:   "Выдох",
			emoji:       "🌬️",
			elapsed:     6,
			total:       6,
			wantCycle:   "Цикл 8/8",
			wantPhase:   "Выдох 🌬️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BreathingPhaseText(
				tt.cycleNum, tt.totalCycles,
				tt.phaseName, tt.emoji,
				tt.elapsed, tt.total,
			)

			if !strings.Contains(result, tt.wantCycle) {
				t.Errorf("expected %q in result, got %q", tt.wantCycle, result)
			}
			if !strings.Contains(result, tt.wantPhase) {
				t.Errorf("expected %q in result, got %q", tt.wantPhase, result)
			}
			// Progress bar should be present
			if !strings.Contains(result, "▓") && !strings.Contains(result, "░") {
				t.Errorf("expected progress bar in result, got %q", result)
			}
		})
	}
}
