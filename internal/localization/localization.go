package localization

import "fmt"

const (
	LangEN      = "en"
	LangRU      = "ru"
	LangDE      = "de"
	LangFR      = "fr"
	LangUK      = "uk"
	LangBE      = "be"
	DefaultLang = LangEN
)

// Messages contains all localized strings.
type Messages struct {
	// Buttons
	Start      string
	Back       string
	Stop       string
	Next       string
	Done       string
	FeelBetter string
	Repeat     string
	BackToMenu string

	// Main menu
	MainMenuText  string
	MenuBreathing string
	MenuGrounding string
	MenuGuided    string
	MenuPMR       string
	MenuThought   string
	MenuInfo      string
	MenuLang      string // "🌐 Language"

	// Breathing
	BreathingIntro      string
	BreathingCompletion string
	BreathingThanks     string
	BreathingCycle      string
	BreathingInhale     string
	BreathingHold       string
	BreathingExhale     string
	BreathingPause      string

	// Grounding
	GroundingIntro      string
	GroundingStep       string
	GroundingCompletion string
	GroundingThanks     string
	GroundingSee        string
	GroundingHear       string
	GroundingFeel       string
	GroundingThings     string
	GroundingSounds     string
	GroundingSensations string

	// Grounding steps (5-4-3-2-1)
	GroundingStep1Title string
	GroundingStep1Desc  string
	GroundingStep2Title string
	GroundingStep2Desc  string
	GroundingStep3Title string
	GroundingStep3Desc  string
	GroundingStep4Title string
	GroundingStep4Desc  string
	GroundingStep5Title string
	GroundingStep5Desc  string

	// Guided breathing
	GuidedIntro      string
	GuidedCompletion string
	GuidedThanks     string
	GuidedStopped    string

	// Breathing patterns
	PatternBoxName        string
	PatternBoxDesc        string
	PatternRelaxingName   string
	PatternRelaxingDesc   string
	PatternEnergizingName string
	PatternEnergizingDesc string
	PatternQuickName      string
	PatternQuickDesc      string

	// PMR
	PMRIntro      string
	PMRTense      string
	PMRRelax      string
	PMRCompletion string
	PMRStopped    string
	PMRThanks     string

	// Thought labeling
	ThoughtIntro      string
	ThoughtPrompt     string
	ThoughtCategories string
	ThoughtResult     string
	ThoughtCompletion string
	ThoughtThanks     string

	// Info
	InfoText string

	// Language selection
	LangSelectTitle string

	// Localized words
	CycleWord string // "Цикл" / "Cycle"
}

// FormatBreathingPhase returns formatted breathing phase text.
func (m Messages) FormatBreathingPhase(cycleNum, totalCycles int, phaseName, emoji string, progress string) string {
	return fmt.Sprintf("*%s %d/%d*\n%s %s\n%s", m.CycleWord, cycleNum, totalCycles, phaseName, emoji, progress)
}

// FormatGuidedPhase returns formatted guided breathing phase text.
func (m Messages) FormatGuidedPhase(patternEmoji, patternName string, cycleNum, totalCycles int, phaseEmoji, phaseName string, progress string) string {
	return fmt.Sprintf("%s *%s*\n*%s %d/%d*\n%s %s\n%s",
		patternEmoji, patternName, m.CycleWord, cycleNum, totalCycles, phaseName, phaseEmoji, progress)
}

// Localizer provides localization services.
type Localizer struct {
	languages          map[string]Messages
	supportedLanguages []string
	defaultLang        string
}

// NewLocalizer creates a new Localizer with all supported languages.
func NewLocalizer() *Localizer {
	return &Localizer{
		languages: map[string]Messages{
			LangEN: messagesEN,
			LangRU: messagesRU,
			LangDE: messagesDE,
			LangFR: messagesFR,
			LangUK: messagesUK,
			LangBE: messagesBE,
		},
		supportedLanguages: []string{LangEN, LangRU, LangDE, LangFR, LangUK, LangBE},
		defaultLang:        DefaultLang,
	}
}

// Get returns messages for the given language code.
func (l *Localizer) Get(langCode string) Messages {
	if len(langCode) >= 2 {
		langCode = langCode[:2]
	}
	if msgs, ok := l.languages[langCode]; ok {
		return msgs
	}
	return l.languages[l.defaultLang]
}

// SupportedLang checks if language is supported and returns normalized code.
func (l *Localizer) SupportedLang(langCode string) string {
	if len(langCode) >= 2 {
		langCode = langCode[:2]
	}
	if _, ok := l.languages[langCode]; ok {
		return langCode
	}
	return l.defaultLang
}

// SupportedLanguages returns list of supported language codes.
func (l *Localizer) SupportedLanguages() []string {
	return l.supportedLanguages
}

// LangFlag returns emoji flag for language.
func (l *Localizer) LangFlag(lang string) string {
	switch lang {
	case LangEN:
		return "🇬🇧"
	case LangRU:
		return "🇷🇺"
	case LangDE:
		return "🇩🇪"
	case LangFR:
		return "🇫🇷"
	case LangUK:
		return "🇺🇦"
	case LangBE:
		return "🇧🇾"
	default:
		return "🌐"
	}
}

// LangName returns native name for language.
func (l *Localizer) LangName(lang string) string {
	switch lang {
	case LangEN:
		return "English"
	case LangRU:
		return "Русский"
	case LangDE:
		return "Deutsch"
	case LangFR:
		return "Français"
	case LangUK:
		return "Українська"
	case LangBE:
		return "Беларуская"
	default:
		return lang
	}
}

// DefaultLocalizer is a global instance for backward compatibility.
// Deprecated: Use NewLocalizer() instead.
var DefaultLocalizer = NewLocalizer()

// SupportedLanguages lists all supported languages in display order.
// Deprecated: Use localizer.SupportedLanguages() instead.
var SupportedLanguages = []string{LangEN, LangRU, LangDE, LangFR, LangUK, LangBE}

// Get returns messages for the given language code.
// Deprecated: Use localizer.Get() instead.
func Get(langCode string) Messages {
	return DefaultLocalizer.Get(langCode)
}

// SupportedLang checks if language is supported and returns normalized code.
// Deprecated: Use localizer.SupportedLang() instead.
func SupportedLang(langCode string) string {
	return DefaultLocalizer.SupportedLang(langCode)
}

// LangFlag returns emoji flag for language.
// Deprecated: Use localizer.LangFlag() instead.
func LangFlag(lang string) string {
	return DefaultLocalizer.LangFlag(lang)
}

// LangName returns native name for language.
// Deprecated: Use localizer.LangName() instead.
func LangName(lang string) string {
	return DefaultLocalizer.LangName(lang)
}
