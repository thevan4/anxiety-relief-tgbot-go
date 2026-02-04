package messages

import (
	"strings"
	"testing"
)

func TestVisualizationIntro(t *testing.T) {
	result := VisualizationIntro("test scenes info")

	if !strings.Contains(result, "Мирная визуализация") {
		t.Errorf("expected title, got %q", result)
	}
	if !strings.Contains(result, "test scenes info") {
		t.Errorf("expected scenes info, got %q", result)
	}
}

func TestVisualizationStepWithProgress(t *testing.T) {
	tests := []struct {
		name        string
		emoji       string
		sceneName   string
		stepNum     int
		totalSteps  int
		instruction string
		elapsed     int
		totalSec    int
	}{
		{
			name:        "first step start",
			emoji:       "🏖️",
			sceneName:   "Пляж",
			stepNum:     1,
			totalSteps:  5,
			instruction: "Представьте тёплый песок",
			elapsed:     0,
			totalSec:    15,
		},
		{
			name:        "middle step",
			emoji:       "🌲",
			sceneName:   "Лес",
			stepNum:     3,
			totalSteps:  5,
			instruction: "Слушайте пение птиц",
			elapsed:     7,
			totalSec:    15,
		},
		{
			name:        "last step end",
			emoji:       "🏔️",
			sceneName:   "Горы",
			stepNum:     5,
			totalSteps:  5,
			instruction: "Вдохните горный воздух",
			elapsed:     15,
			totalSec:    15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VisualizationStepWithProgress(
				tt.emoji, tt.sceneName,
				tt.stepNum, tt.totalSteps,
				tt.instruction, tt.elapsed, tt.totalSec,
			)

			if !strings.Contains(result, tt.emoji) {
				t.Errorf("expected emoji %q, got %q", tt.emoji, result)
			}
			if !strings.Contains(result, tt.sceneName) {
				t.Errorf("expected scene name %q, got %q", tt.sceneName, result)
			}
			if !strings.Contains(result, tt.instruction) {
				t.Errorf("expected instruction %q, got %q", tt.instruction, result)
			}
			// Timer countdown should be present
			if !strings.Contains(result, "⏱️") && !strings.Contains(result, "✨") {
				t.Errorf("expected timer countdown in result, got %q", result)
			}
		})
	}
}

func TestVisualizationIntroWithProgress(t *testing.T) {
	result := VisualizationIntroWithProgress(
		"🏖️", "Пляж", "Тёплый песчаный пляж", "Шум волн",
		3, 5,
	)

	if !strings.Contains(result, "🏖️") {
		t.Errorf("expected emoji, got %q", result)
	}
	if !strings.Contains(result, "Пляж") {
		t.Errorf("expected scene name, got %q", result)
	}
	if !strings.Contains(result, "Тёплый песчаный пляж") {
		t.Errorf("expected description, got %q", result)
	}
	if !strings.Contains(result, "Шум волн") {
		t.Errorf("expected atmosphere, got %q", result)
	}
	// Timer countdown should be present
	if !strings.Contains(result, "⏱️") && !strings.Contains(result, "✨") {
		t.Errorf("expected timer countdown in result, got %q", result)
	}
}

func TestVisualizationCompletion(t *testing.T) {
	result := VisualizationCompletion("Лесная поляна")

	if !strings.Contains(result, "Лесная поляна") {
		t.Errorf("expected scene name, got %q", result)
	}
	if !strings.Contains(result, "Отлично") {
		t.Errorf("expected completion message, got %q", result)
	}
}
