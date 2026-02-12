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
				{"Pause", msgs.Pause},
				{"Resume", msgs.Resume},
				{"PauseText", msgs.PauseText},
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

func TestFormatBreathingPhase(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		lang := lang
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			m := l.Get(lang)
			result := m.FormatBreathingPhase(1, 8, "Inhale", "🫁", "████░░░░")

			assertContains(t, result, "1", "8", "Inhale", "🫁", "████░░░░")
			assertContains(t, result, m.CycleWord)
		})
	}
}

func TestFormatGuidedPhase(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		lang := lang
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			m := l.Get(lang)
			result := m.FormatGuidedPhase("📦", "Box Breathing", 2, 5, "🫁", "Inhale", "████░░░░")

			assertContains(t, result, "📦", "Box Breathing", "2", "5", "🫁", "Inhale")
			assertContains(t, result, m.CycleWord)
		})
	}
}

func TestPMRMessagesNotEmpty(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		lang := lang
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			m := l.Get(lang)

			fields := []struct {
				name  string
				value string
			}{
				{"PMRIntro", m.PMRIntro},
				{"PMRTense", m.PMRTense},
				{"PMRRelax", m.PMRRelax},
				{"PMRCompletion", m.PMRCompletion},
				{"PMRHandsName", m.PMRHandsName},
				{"PMRHandsTense", m.PMRHandsTense},
				{"PMRHandsRelax", m.PMRHandsRelax},
			}
			for _, f := range fields {
				if f.value == "" {
					t.Errorf("%s: %s is empty", lang, f.name)
				}
			}
		})
	}
}

func TestGuidedBreathingMessagesNotEmpty(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		lang := lang
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			m := l.Get(lang)

			fields := []struct {
				name  string
				value string
			}{
				{"GuidedIntro", m.GuidedIntro},
				{"GuidedCompletion", m.GuidedCompletion},
				{"PatternBoxName", m.PatternBoxName},
				{"PatternBoxDesc", m.PatternBoxDesc},
				{"PatternRelaxingName", m.PatternRelaxingName},
				{"PatternQuickName", m.PatternQuickName},
				{"GuidedInhale", m.GuidedInhale},
				{"GuidedExhale", m.GuidedExhale},
			}
			for _, f := range fields {
				if f.value == "" {
					t.Errorf("%s: %s is empty", lang, f.name)
				}
			}
		})
	}
}

func TestGroundingMessagesNotEmpty(t *testing.T) {
	t.Parallel()
	l := NewLocalizer()

	for _, lang := range l.SupportedLanguages() {
		lang := lang
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			m := l.Get(lang)

			fields := []struct {
				name  string
				value string
			}{
				{"GroundingStep1Title", m.GroundingStep1Title},
				{"GroundingStep1Desc", m.GroundingStep1Desc},
				{"GroundingStep2Title", m.GroundingStep2Title},
				{"GroundingStep3Title", m.GroundingStep3Title},
				{"GroundingStep4Title", m.GroundingStep4Title},
				{"GroundingStep5Title", m.GroundingStep5Title},
			}
			for _, f := range fields {
				if f.value == "" {
					t.Errorf("%s: %s is empty", lang, f.name)
				}
			}
		})
	}
}
