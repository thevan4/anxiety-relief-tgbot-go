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

	// Visualization
	MenuVisualization       string
	VisualizationIntro      string
	VisualizationStopped    string
	VisualizationCompletion string
	VisualizationThanks     string
	VisualizationAtmosphere string
	VisualizationCloseEyes  string
	VisualizationStepFmt    string // "Шаг %d/%d"

	// Visualization scenes
	SceneMountainName  string
	SceneMountainDesc  string
	SceneMountainAtmo  string
	SceneMountainStep1 string
	SceneMountainStep2 string
	SceneMountainStep3 string
	SceneMountainStep4 string
	SceneMountainStep5 string
	SceneMountainStep6 string
	SceneMountainStep7 string

	SceneForestName  string
	SceneForestDesc  string
	SceneForestAtmo  string
	SceneForestStep1 string
	SceneForestStep2 string
	SceneForestStep3 string
	SceneForestStep4 string
	SceneForestStep5 string
	SceneForestStep6 string
	SceneForestStep7 string

	SceneBeachName  string
	SceneBeachDesc  string
	SceneBeachAtmo  string
	SceneBeachStep1 string
	SceneBeachStep2 string
	SceneBeachStep3 string
	SceneBeachStep4 string
	SceneBeachStep5 string
	SceneBeachStep6 string
	SceneBeachStep7 string

	SceneGardenName  string
	SceneGardenDesc  string
	SceneGardenAtmo  string
	SceneGardenStep1 string
	SceneGardenStep2 string
	SceneGardenStep3 string
	SceneGardenStep4 string
	SceneGardenStep5 string
	SceneGardenStep6 string
	SceneGardenStep7 string

	SceneStarryName  string
	SceneStarryDesc  string
	SceneStarryAtmo  string
	SceneStarryStep1 string
	SceneStarryStep2 string
	SceneStarryStep3 string
	SceneStarryStep4 string
	SceneStarryStep5 string
	SceneStarryStep6 string
	SceneStarryStep7 string

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

// FormatVisualizationSceneIntro returns formatted scene intro with timer.
func (m Messages) FormatVisualizationSceneIntro(emoji, name, description, atmosphere, timer string) string {
	return fmt.Sprintf("%s *%s*\n\n_%s_\n\n🌿 *%s:* %s\n\n%s\n\n%s",
		emoji, name, description, m.VisualizationAtmosphere, atmosphere, m.VisualizationCloseEyes, timer)
}

// FormatVisualizationStep returns formatted visualization step with timer.
func (m Messages) FormatVisualizationStep(emoji, name string, stepNum, totalSteps int, instruction, timer string) string {
	return fmt.Sprintf("%s *%s*\n\n*"+m.VisualizationStepFmt+"*\n\n%s\n\n%s",
		emoji, name, stepNum, totalSteps, instruction, timer)
}

// FormatVisualizationCompletion returns formatted completion message.
func (m Messages) FormatVisualizationCompletion(sceneName string) string {
	return fmt.Sprintf(m.VisualizationCompletion, sceneName)
}

// VisualizationScene represents a localized visualization scene.
type VisualizationScene struct {
	ID          string
	Name        string
	Emoji       string
	Description string
	Atmosphere  string
	Steps       []VisualizationStep
}

// VisualizationStep represents one step in visualization.
type VisualizationStep struct {
	Number      int
	Instruction string
}

// GetVisualizationScenes returns localized visualization scenes.
func (m Messages) GetVisualizationScenes() []VisualizationScene {
	return []VisualizationScene{
		{
			ID: "mountain", Name: m.SceneMountainName, Emoji: "🏔️",
			Description: m.SceneMountainDesc, Atmosphere: m.SceneMountainAtmo,
			Steps: []VisualizationStep{
				{1, m.SceneMountainStep1}, {2, m.SceneMountainStep2}, {3, m.SceneMountainStep3},
				{4, m.SceneMountainStep4}, {5, m.SceneMountainStep5}, {6, m.SceneMountainStep6},
				{7, m.SceneMountainStep7},
			},
		},
		{
			ID: "forest", Name: m.SceneForestName, Emoji: "🌲",
			Description: m.SceneForestDesc, Atmosphere: m.SceneForestAtmo,
			Steps: []VisualizationStep{
				{1, m.SceneForestStep1}, {2, m.SceneForestStep2}, {3, m.SceneForestStep3},
				{4, m.SceneForestStep4}, {5, m.SceneForestStep5}, {6, m.SceneForestStep6},
				{7, m.SceneForestStep7},
			},
		},
		{
			ID: "beach", Name: m.SceneBeachName, Emoji: "🏖️",
			Description: m.SceneBeachDesc, Atmosphere: m.SceneBeachAtmo,
			Steps: []VisualizationStep{
				{1, m.SceneBeachStep1}, {2, m.SceneBeachStep2}, {3, m.SceneBeachStep3},
				{4, m.SceneBeachStep4}, {5, m.SceneBeachStep5}, {6, m.SceneBeachStep6},
				{7, m.SceneBeachStep7},
			},
		},
		{
			ID: "garden", Name: m.SceneGardenName, Emoji: "🌸",
			Description: m.SceneGardenDesc, Atmosphere: m.SceneGardenAtmo,
			Steps: []VisualizationStep{
				{1, m.SceneGardenStep1}, {2, m.SceneGardenStep2}, {3, m.SceneGardenStep3},
				{4, m.SceneGardenStep4}, {5, m.SceneGardenStep5}, {6, m.SceneGardenStep6},
				{7, m.SceneGardenStep7},
			},
		},
		{
			ID: "starry", Name: m.SceneStarryName, Emoji: "🌌",
			Description: m.SceneStarryDesc, Atmosphere: m.SceneStarryAtmo,
			Steps: []VisualizationStep{
				{1, m.SceneStarryStep1}, {2, m.SceneStarryStep2}, {3, m.SceneStarryStep3},
				{4, m.SceneStarryStep4}, {5, m.SceneStarryStep5}, {6, m.SceneStarryStep6},
				{7, m.SceneStarryStep7},
			},
		},
	}
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
//
// Deprecated: Use NewLocalizer() instead.
var DefaultLocalizer = NewLocalizer()

// SupportedLanguages lists all supported languages in display order.
//
// Deprecated: Use localizer.SupportedLanguages() instead.
var SupportedLanguages = []string{LangEN, LangRU, LangDE, LangFR, LangUK, LangBE}

// Get returns messages for the given language code.
//
// Deprecated: Use localizer.Get() instead.
func Get(langCode string) Messages {
	return DefaultLocalizer.Get(langCode)
}

// SupportedLang checks if language is supported and returns normalized code.
//
// Deprecated: Use localizer.SupportedLang() instead.
func SupportedLang(langCode string) string {
	return DefaultLocalizer.SupportedLang(langCode)
}

// LangFlag returns emoji flag for language.
//
// Deprecated: Use localizer.LangFlag() instead.
func LangFlag(lang string) string {
	return DefaultLocalizer.LangFlag(lang)
}

// LangName returns native name for language.
//
// Deprecated: Use localizer.LangName() instead.
func LangName(lang string) string {
	return DefaultLocalizer.LangName(lang)
}
