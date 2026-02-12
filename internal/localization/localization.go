// Package localization provides multi-language message support for the bot.
package localization

import "fmt"

const (
	// LangEN is the English language code.
	LangEN = "en"
	// LangRU is the Russian language code.
	LangRU = "ru"
	// LangDE is the German language code.
	LangDE = "de"
	// LangFR is the French language code.
	LangFR = "fr"
	// LangUK is the Ukrainian language code.
	LangUK = "uk"
	// LangBE is the Belarusian language code.
	LangBE = "be"
	// DefaultLang is the default language if user hasn't selected one.
	DefaultLang = LangEN

	langCodeMinLen = 2
)

// Messages contains all localized strings.
type Messages struct {
	// Buttons
	Start  string // "Начать" — starts new session from holder
	Back   string
	Stop   string
	Next   string
	Done   string
	Repeat string
	Pause  string // "⏸ Pause" — pauses running exercise
	Resume string // "▶️ Continue" — resumes paused exercise

	// Pause screen
	PauseText string

	// Welcome holder (after /start) — static message with bot info
	HolderText string

	// Session expired message
	SessionExpired string

	// Main menu
	MainMenuText  string
	MenuBreathing string
	MenuGrounding string
	MenuGuided    string
	MenuPMR       string
	MenuLang      string // "🌐 Language"

	// Breathing
	BreathingIntro      string
	BreathingCompletion string
	BreathingInhale     string
	BreathingHold       string
	BreathingExhale     string

	// Grounding
	GroundingIntro      string
	GroundingCompletion string

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

	// Guided breathing — pattern intros
	PatternBoxIntro        string
	PatternRelaxingIntro   string
	PatternEnergizingIntro string
	PatternQuickIntro      string

	// Breathing patterns
	PatternBoxName        string
	PatternBoxDesc        string
	PatternRelaxingName   string
	PatternRelaxingDesc   string
	PatternEnergizingName string
	PatternEnergizingDesc string
	PatternQuickName      string
	PatternQuickDesc      string

	// Guided breathing phases
	GuidedInhale  string
	GuidedHoldIn  string
	GuidedExhale  string
	GuidedHoldOut string

	// PMR
	PMRIntro      string
	PMRTense      string
	PMRRelax      string
	PMRCompletion string

	// PMR muscle groups
	PMRHandsName  string
	PMRHandsTense string
	PMRHandsRelax string

	PMRForearmsName  string
	PMRForearmsTense string
	PMRForearmsRelax string

	PMRForeheadName  string
	PMRForeheadTense string
	PMRForeheadRelax string

	PMREyesName  string
	PMREyesTense string
	PMREyesRelax string

	PMRJawName  string
	PMRJawTense string
	PMRJawRelax string

	PMRNeckName  string
	PMRNeckTense string
	PMRNeckRelax string

	PMRChestName  string
	PMRChestTense string
	PMRChestRelax string

	PMRStomachName  string
	PMRStomachTense string
	PMRStomachRelax string

	PMRThighsName  string
	PMRThighsTense string
	PMRThighsRelax string

	PMRCalvesName  string
	PMRCalvesTense string
	PMRCalvesRelax string

	// Language selection
	LangSelectTitle string

	// Localized words
	CycleWord string // "Цикл" / "Cycle"
}

// FormatBreathingPhase returns formatted breathing phase text.
func (m Messages) FormatBreathingPhase(
	cycleNum, totalCycles int,
	phaseName, emoji, progress string,
) string {
	return fmt.Sprintf(
		"*%s %d/%d*\n%s %s\n%s",
		m.CycleWord,
		cycleNum,
		totalCycles,
		phaseName,
		emoji,
		progress,
	)
}

// FormatGuidedPhase returns formatted guided breathing phase text.
func (m Messages) FormatGuidedPhase(
	patternEmoji, patternName string,
	cycleNum, totalCycles int,
	phaseEmoji, phaseName, progress string,
) string {
	return fmt.Sprintf(
		"%s *%s*\n*%s %d/%d*\n%s %s\n%s",
		patternEmoji,
		patternName,
		m.CycleWord,
		cycleNum,
		totalCycles,
		phaseName,
		phaseEmoji,
		progress,
	)
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
	if len(langCode) >= langCodeMinLen {
		langCode = langCode[:langCodeMinLen]
	}
	if msgs, ok := l.languages[langCode]; ok {
		return msgs
	}
	return l.languages[l.defaultLang]
}

// SupportedLang checks if language is supported and returns normalized code.
func (l *Localizer) SupportedLang(langCode string) string {
	if len(langCode) >= langCodeMinLen {
		langCode = langCode[:langCodeMinLen]
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
