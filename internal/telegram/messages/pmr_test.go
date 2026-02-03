package messages

import (
	"strings"
	"testing"
)

func TestPMRIntro(t *testing.T) {
	result := PMRIntro(10)

	if !strings.Contains(result, "Прогрессивная мышечная релаксация") {
		t.Errorf("expected title, got %q", result)
	}
	if !strings.Contains(result, "10 групп мышц") {
		t.Errorf("expected muscle groups count, got %q", result)
	}
}

func TestPMRTensePhaseWithProgress(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		elapsed     int
		totalSec    int
	}{
		{
			name:        "start",
			instruction: "Сожмите кулаки",
			elapsed:     0,
			totalSec:    7,
		},
		{
			name:        "middle",
			instruction: "Напрягите бицепсы",
			elapsed:     3,
			totalSec:    7,
		},
		{
			name:        "end",
			instruction: "Поднимите брови",
			elapsed:     7,
			totalSec:    7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PMRTensePhaseWithProgress(tt.instruction, tt.elapsed, tt.totalSec)

			if !strings.Contains(result, "НАПРЯГИТЕ") {
				t.Errorf("expected tense label, got %q", result)
			}
			if !strings.Contains(result, tt.instruction) {
				t.Errorf("expected instruction %q, got %q", tt.instruction, result)
			}
			// Progress bar should be present
			if !strings.Contains(result, "▓") && !strings.Contains(result, "░") {
				t.Errorf("expected progress bar in result, got %q", result)
			}
		})
	}
}

func TestPMRRelaxPhaseWithProgress(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		elapsed     int
		totalSec    int
	}{
		{
			name:        "start",
			instruction: "Разожмите кулаки",
			elapsed:     0,
			totalSec:    15,
		},
		{
			name:        "middle",
			instruction: "Расслабьте руки",
			elapsed:     7,
			totalSec:    15,
		},
		{
			name:        "end",
			instruction: "Опустите плечи",
			elapsed:     15,
			totalSec:    15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PMRRelaxPhaseWithProgress(tt.instruction, tt.elapsed, tt.totalSec)

			if !strings.Contains(result, "РАССЛАБЬТЕ") {
				t.Errorf("expected relax label, got %q", result)
			}
			if !strings.Contains(result, tt.instruction) {
				t.Errorf("expected instruction %q, got %q", tt.instruction, result)
			}
			// Progress bar should be present
			if !strings.Contains(result, "▓") && !strings.Contains(result, "░") {
				t.Errorf("expected progress bar in result, got %q", result)
			}
		})
	}
}

func TestPMRMuscleStep(t *testing.T) {
	result := PMRMuscleStep("✊", "Кисти рук", 1, 10, "test phase text")

	if !strings.Contains(result, "✊") {
		t.Errorf("expected emoji, got %q", result)
	}
	if !strings.Contains(result, "Кисти рук") {
		t.Errorf("expected muscle name, got %q", result)
	}
	if !strings.Contains(result, "1/10") {
		t.Errorf("expected progress indicator, got %q", result)
	}
	if !strings.Contains(result, "test phase text") {
		t.Errorf("expected phase text, got %q", result)
	}
}
