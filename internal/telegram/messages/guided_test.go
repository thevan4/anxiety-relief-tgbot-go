package messages

import (
	"strings"
	"testing"
)

func TestGuidedBreathingIntro(t *testing.T) {
	result := GuidedBreathingIntro("test patterns info")

	if !strings.Contains(result, "Управляемое дыхание") {
		t.Errorf("expected title, got %q", result)
	}
	if !strings.Contains(result, "test patterns info") {
		t.Errorf("expected patterns info, got %q", result)
	}
}

func TestGuidedBreathingPhase(t *testing.T) {
	tests := []struct {
		name         string
		patternEmoji string
		patternName  string
		cycle        int
		totalCycles  int
		phaseEmoji   string
		phaseName    string
		elapsed      int
		totalSec     int
		wantCycle    string
		wantPhase    string
	}{
		{
			name:         "inhale start",
			patternEmoji: "🌊",
			patternName:  "4-7-8",
			cycle:        1,
			totalCycles:  4,
			phaseEmoji:   "🫁",
			phaseName:    "Вдох",
			elapsed:      0,
			totalSec:     4,
			wantCycle:    "Цикл 1/4",
			wantPhase:    "Вдох 🫁",
		},
		{
			name:         "hold middle",
			patternEmoji: "⬜",
			patternName:  "Коробка",
			cycle:        2,
			totalCycles:  4,
			phaseEmoji:   "⏸️",
			phaseName:    "Задержка",
			elapsed:      3,
			totalSec:     7,
			wantCycle:    "Цикл 2/4",
			wantPhase:    "Задержка ⏸️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GuidedBreathingPhase(
				tt.patternEmoji, tt.patternName,
				tt.cycle, tt.totalCycles,
				tt.phaseEmoji, tt.phaseName,
				tt.elapsed, tt.totalSec,
			)

			if !strings.Contains(result, tt.patternEmoji) {
				t.Errorf("expected pattern emoji %q in result, got %q", tt.patternEmoji, result)
			}
			if !strings.Contains(result, tt.patternName) {
				t.Errorf("expected pattern name %q in result, got %q", tt.patternName, result)
			}
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

func TestGuidedBreathingCompletion(t *testing.T) {
	result := GuidedBreathingCompletion("4-7-8")

	if !strings.Contains(result, "4-7-8") {
		t.Errorf("expected pattern name, got %q", result)
	}
	if !strings.Contains(result, "Отлично") {
		t.Errorf("expected completion message, got %q", result)
	}
}
