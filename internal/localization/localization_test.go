package localization

import (
	"strings"
	"testing"
)

func TestLocalizerGet(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := l.SupportedLang(tt.langCode)
			if got != tt.want {
				t.Errorf("SupportedLang(%q) = %q, want %q", tt.langCode, got, tt.want)
			}
		})
	}
}

func TestLocalizerSupportedLanguages(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
		tt := tt
		t.Run(tt.lang, func(t *testing.T) {
			t.Parallel()
			got := l.LangFlag(tt.lang)
			if got != tt.want {
				t.Errorf("LangFlag(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

func TestLocalizerLangName(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.lang, func(t *testing.T) {
			t.Parallel()
			got := l.LangName(tt.lang)
			if got != tt.want {
				t.Errorf("LangName(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

func TestMessagesNotEmpty(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			msgs := l.Get(lang)

			fields := []struct {
				name  string
				value string
			}{
				{"Start", msgs.Start},
				{"Back", msgs.Back},
				{"Stop", msgs.Stop},
				{"Next", msgs.Next},
				{"Done", msgs.Done},
				{"MainMenuText", msgs.MainMenuText},
				{"MenuBreathing", msgs.MenuBreathing},
				{"BreathingIntro", msgs.BreathingIntro},
				{"BreathingInhale", msgs.BreathingInhale},
				{"BreathingHold", msgs.BreathingHold},
				{"BreathingExhale", msgs.BreathingExhale},
				{"GroundingIntro", msgs.GroundingIntro},
			}
			for _, field := range fields {
				if field.value == "" {
					t.Errorf("%s: %s is empty", lang, field.name)
				}
			}
		})
	}
}

func TestVisualizationMessagesNotEmpty(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			msgs := l.Get(lang)

			fields := []struct {
				name  string
				value string
			}{
				{"MenuVisualization", msgs.MenuVisualization},
				{"VisualizationIntro", msgs.VisualizationIntro},
				{"VisualizationStopped", msgs.VisualizationStopped},
				{"VisualizationCompletion", msgs.VisualizationCompletion},
				{"VisualizationThanks", msgs.VisualizationThanks},
				{"VisualizationAtmosphere", msgs.VisualizationAtmosphere},
				{"VisualizationCloseEyes", msgs.VisualizationCloseEyes},
				{"VisualizationStepFmt", msgs.VisualizationStepFmt},
			}
			for _, field := range fields {
				if field.value == "" {
					t.Errorf("%s: %s is empty", lang, field.name)
				}
			}
		})
	}
}

func validateScene(t *testing.T, lang string, scene VisualizationScene, expectedID string) {
	t.Helper()
	if scene.ID != expectedID {
		t.Errorf("%s: scene ID = %q, want %q", lang, scene.ID, expectedID)
	}
	sceneFields := []struct {
		name  string
		value string
	}{
		{"Name", scene.Name},
		{"Description", scene.Description},
		{"Atmosphere", scene.Atmosphere},
	}
	for _, field := range sceneFields {
		if field.value == "" {
			t.Errorf("%s: scene %s %s is empty", lang, scene.ID, field.name)
		}
	}
	if len(scene.Steps) != 7 {
		t.Errorf("%s: scene %s has %d steps, want 7", lang, scene.ID, len(scene.Steps))
	}
	for i, step := range scene.Steps {
		if step.Instruction == "" {
			t.Errorf("%s: scene %s step %d Instruction is empty", lang, scene.ID, i+1)
		}
	}
}

func TestGetVisualizationScenes(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		lang := lang
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			m := l.Get(lang)
			scenes := m.GetVisualizationScenes()

			if len(scenes) != 5 {
				t.Errorf("%s: expected 5 scenes, got %d", lang, len(scenes))
			}

			expectedIDs := []string{"mountain", "forest", "beach", "garden", "starry"}
			for i, scene := range scenes {
				validateScene(t, lang, scene, expectedIDs[i])
			}
		})
	}
}

func TestVisualizationFormatMethods(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()
	m := l.Get("en")

	t.Run("FormatVisualizationSceneIntro", func(t *testing.T) {
		t.Parallel()
		result := m.FormatVisualizationSceneIntro("🏖️", "Beach", "Warm beach", "Sound of waves", "⏱️ 3")
		assertContains(t, result, "🏖️", "Beach", "Warm beach")
	})

	t.Run("FormatVisualizationStep", func(t *testing.T) {
		t.Parallel()
		result := m.FormatVisualizationStep("🏖️", "Beach", 1, 5, "Close your eyes", "⏱️ 10")
		assertContains(t, result, "🏖️", "Beach", "Close your eyes")
	})

	t.Run("FormatVisualizationCompletion", func(t *testing.T) {
		t.Parallel()
		result := m.FormatVisualizationCompletion("Beach Scene")
		assertContains(t, result, "Beach Scene")
	})
}

// assertContains checks that result contains all expected substrings.
func assertContains(t *testing.T, result string, expected ...string) {
	t.Helper()
	if result == "" {
		t.Error("result is empty")
		return
	}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("result missing %q: %q", exp, result)
		}
	}
}

// Test backward compatibility functions.
func TestBackwardCompatibility(t *testing.T) {
	t.Parallel()
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
