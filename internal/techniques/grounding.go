package techniques

const (
	groundingSeeCount   = 5
	groundingTouchCount = 4
	groundingHearCount  = 3
	groundingSmellCount = 2
	groundingTasteCount = 1

	groundingStep1 = 1
	groundingStep2 = 2
	groundingStep3 = 3
	groundingStep4 = 4
	groundingStep5 = 5
)

// GetGroundingSteps returns the 5-4-3-2-1 grounding technique steps.
func GetGroundingSteps() []GroundingStep {
	return []GroundingStep{
		{Number: groundingStep1, Sense: "sight", Count: groundingSeeCount, Emoji: "👁️"},
		{Number: groundingStep2, Sense: "touch", Count: groundingTouchCount, Emoji: "🤚"},
		{Number: groundingStep3, Sense: "hearing", Count: groundingHearCount, Emoji: "👂"},
		{Number: groundingStep4, Sense: "smell", Count: groundingSmellCount, Emoji: "👃"},
		{Number: groundingStep5, Sense: "taste", Count: groundingTasteCount, Emoji: "👅"},
	}
}
