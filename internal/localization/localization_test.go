package localization

import (
	"testing"
)

func TestLocalizerGet(t *testing.T) {
	l := NewLocalizer()

	tests := []struct {
		name     string
		langCode string
		wantLang string
	}{
		{"english", "en", "en"},
		{"russian", "ru", "ru"},
		{"german", "de", "de"},
		{"french", "fr", "fr"},
		{"ukrainian", "uk", "uk"},
		{"belarusian", "be", "be"},
		{"english with region", "en-US", "en"},
		{"unknown defaults to english", "xx", "en"},
		{"empty defaults to english", "", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgs := l.Get(tt.langCode)
			// Check that we got non-empty messages
			if msgs.MainMenuText == "" {
				t.Errorf("Get(%q) returned empty MainMenuText", tt.langCode)
			}
			if msgs.Start == "" {
				t.Errorf("Get(%q) returned empty Start button", tt.langCode)
			}
		})
	}
}

func TestLocalizerSupportedLang(t *testing.T) {
	l := NewLocalizer()

	tests := []struct {
		name     string
		langCode string
		want     string
	}{
		{"english", "en", "en"},
		{"russian", "ru", "ru"},
		{"german", "de", "de"},
		{"french", "fr", "fr"},
		{"ukrainian", "uk", "uk"},
		{"belarusian", "be", "be"},
		{"english with region", "en-US", "en"},
		{"russian with region", "ru-RU", "ru"},
		{"unknown defaults to english", "xx", "en"},
		{"empty defaults to english", "", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := l.SupportedLang(tt.langCode)
			if got != tt.want {
				t.Errorf("SupportedLang(%q) = %q, want %q", tt.langCode, got, tt.want)
			}
		})
	}
}

func TestLocalizerSupportedLanguages(t *testing.T) {
	l := NewLocalizer()
	langs := l.SupportedLanguages()

	if len(langs) != 6 {
		t.Errorf("SupportedLanguages() returned %d languages, want 6", len(langs))
	}

	expected := []string{LangEN, LangRU, LangDE, LangFR, LangUK, LangBE}
	for i, lang := range expected {
		if langs[i] != lang {
			t.Errorf("SupportedLanguages()[%d] = %q, want %q", i, langs[i], lang)
		}
	}
}

func TestLocalizerLangFlag(t *testing.T) {
	l := NewLocalizer()

	tests := []struct {
		lang string
		want string
	}{
		{LangEN, "🇬🇧"},
		{LangRU, "🇷🇺"},
		{LangDE, "🇩🇪"},
		{LangFR, "🇫🇷"},
		{LangUK, "🇺🇦"},
		{LangBE, "🇧🇾"},
		{"xx", "🌐"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := l.LangFlag(tt.lang)
			if got != tt.want {
				t.Errorf("LangFlag(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

func TestLocalizerLangName(t *testing.T) {
	l := NewLocalizer()

	tests := []struct {
		lang string
		want string
	}{
		{LangEN, "English"},
		{LangRU, "Русский"},
		{LangDE, "Deutsch"},
		{LangFR, "Français"},
		{LangUK, "Українська"},
		{LangBE, "Беларуская"},
		{"xx", "xx"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := l.LangName(tt.lang)
			if got != tt.want {
				t.Errorf("LangName(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

func TestMessagesNotEmpty(t *testing.T) {
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		t.Run(lang, func(t *testing.T) {
			msgs := l.Get(lang)

			// Check all button texts are non-empty
			if msgs.Start == "" {
				t.Errorf("%s: Start is empty", lang)
			}
			if msgs.Back == "" {
				t.Errorf("%s: Back is empty", lang)
			}
			if msgs.Stop == "" {
				t.Errorf("%s: Stop is empty", lang)
			}
			if msgs.Next == "" {
				t.Errorf("%s: Next is empty", lang)
			}
			if msgs.Done == "" {
				t.Errorf("%s: Done is empty", lang)
			}

			// Check main menu
			if msgs.MainMenuText == "" {
				t.Errorf("%s: MainMenuText is empty", lang)
			}
			if msgs.MenuBreathing == "" {
				t.Errorf("%s: MenuBreathing is empty", lang)
			}

			// Check breathing
			if msgs.BreathingIntro == "" {
				t.Errorf("%s: BreathingIntro is empty", lang)
			}
			if msgs.BreathingInhale == "" {
				t.Errorf("%s: BreathingInhale is empty", lang)
			}
			if msgs.BreathingHold == "" {
				t.Errorf("%s: BreathingHold is empty", lang)
			}
			if msgs.BreathingExhale == "" {
				t.Errorf("%s: BreathingExhale is empty", lang)
			}

			// Check grounding
			if msgs.GroundingIntro == "" {
				t.Errorf("%s: GroundingIntro is empty", lang)
			}
		})
	}
}

func TestVisualizationMessagesNotEmpty(t *testing.T) {
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		t.Run(lang, func(t *testing.T) {
			msgs := l.Get(lang)

			if msgs.MenuVisualization == "" {
				t.Errorf("%s: MenuVisualization is empty", lang)
			}
			if msgs.VisualizationIntro == "" {
				t.Errorf("%s: VisualizationIntro is empty", lang)
			}
			if msgs.VisualizationStopped == "" {
				t.Errorf("%s: VisualizationStopped is empty", lang)
			}
			if msgs.VisualizationCompletion == "" {
				t.Errorf("%s: VisualizationCompletion is empty", lang)
			}
			if msgs.VisualizationThanks == "" {
				t.Errorf("%s: VisualizationThanks is empty", lang)
			}
			if msgs.VisualizationAtmosphere == "" {
				t.Errorf("%s: VisualizationAtmosphere is empty", lang)
			}
			if msgs.VisualizationCloseEyes == "" {
				t.Errorf("%s: VisualizationCloseEyes is empty", lang)
			}
			if msgs.VisualizationStepFmt == "" {
				t.Errorf("%s: VisualizationStepFmt is empty", lang)
			}
		})
	}
}

func TestGetVisualizationScenes(t *testing.T) {
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		t.Run(lang, func(t *testing.T) {
			m := l.Get(lang)
			scenes := m.GetVisualizationScenes()

			if len(scenes) != 5 {
				t.Errorf("%s: expected 5 scenes, got %d", lang, len(scenes))
			}

			expectedIDs := []string{"mountain", "forest", "beach", "garden", "starry"}
			for i, scene := range scenes {
				if scene.ID != expectedIDs[i] {
					t.Errorf("%s: scene %d ID = %q, want %q", lang, i, scene.ID, expectedIDs[i])
				}
				if scene.Name == "" {
					t.Errorf("%s: scene %s Name is empty", lang, scene.ID)
				}
				if scene.Description == "" {
					t.Errorf("%s: scene %s Description is empty", lang, scene.ID)
				}
				if scene.Atmosphere == "" {
					t.Errorf("%s: scene %s Atmosphere is empty", lang, scene.ID)
				}
				if len(scene.Steps) != 7 {
					t.Errorf("%s: scene %s has %d steps, want 7", lang, scene.ID, len(scene.Steps))
				}
				for j, step := range scene.Steps {
					if step.Instruction == "" {
						t.Errorf("%s: scene %s step %d Instruction is empty", lang, scene.ID, j+1)
					}
				}
			}
		})
	}
}

func TestVisualizationFormatMethods(t *testing.T) {
	l := NewLocalizer()
	m := l.Get("en")

	t.Run("FormatVisualizationSceneIntro", func(t *testing.T) {
		result := m.FormatVisualizationSceneIntro("🏖️", "Beach", "Warm beach", "Sound of waves", "⏱️ 3")
		if result == "" {
			t.Error("FormatVisualizationSceneIntro returned empty string")
		}
		if !contains(result, "🏖️") || !contains(result, "Beach") || !contains(result, "Warm beach") {
			t.Errorf("FormatVisualizationSceneIntro missing expected content: %q", result)
		}
	})

	t.Run("FormatVisualizationStep", func(t *testing.T) {
		result := m.FormatVisualizationStep("🏖️", "Beach", 1, 5, "Close your eyes", "⏱️ 10")
		if result == "" {
			t.Error("FormatVisualizationStep returned empty string")
		}
		if !contains(result, "🏖️") || !contains(result, "Beach") || !contains(result, "Close your eyes") {
			t.Errorf("FormatVisualizationStep missing expected content: %q", result)
		}
	})

	t.Run("FormatVisualizationCompletion", func(t *testing.T) {
		result := m.FormatVisualizationCompletion("Beach Scene")
		if result == "" {
			t.Error("FormatVisualizationCompletion returned empty string")
		}
		if !contains(result, "Beach Scene") {
			t.Errorf("FormatVisualizationCompletion missing scene name: %q", result)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Test backward compatibility functions.
func TestBackwardCompatibility(t *testing.T) {
	// Test global Get function
	msgs := Get("en")
	if msgs.Start == "" {
		t.Error("Get() returned empty Start")
	}

	// Test global SupportedLang function
	if SupportedLang("ru") != "ru" {
		t.Error("SupportedLang() failed")
	}

	// Test global LangFlag function
	if LangFlag("en") != "🇬🇧" {
		t.Error("LangFlag() failed")
	}

	// Test global LangName function
	if LangName("en") != "English" {
		t.Error("LangName() failed")
	}

	// Test SupportedLanguages var
	if len(SupportedLanguages) != 6 {
		t.Errorf("SupportedLanguages has %d elements, want 6", len(SupportedLanguages))
	}
}
