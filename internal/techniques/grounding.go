package techniques

const (
	groundingSeeCount   = 5
	groundingTouchCount = 4
	groundingHearCount  = 3
	groundingSmellCount = 2
	groundingTasteCount = 1
)

// GetGroundingSteps returns the 5-4-3-2-1 grounding technique steps.
func GetGroundingSteps() []GroundingStep {
	touchExample := "Назовите 4 вещи, которые вы ощущаете телом.\n\n" +
		"Например: ноги на полу, спина на стуле, руки на столе, одежда на теле."
	return []GroundingStep{
		{
			Number:      groundingSeeCount,
			Title:       "Шаг 1: Зрение",
			Description: "Назовите 5 вещей, которые вы видите вокруг себя.\n\nНапример: стол, лампа, книга, окно, стул.",
			Emoji:       "👁️",
		},
		{
			Number:      groundingTouchCount,
			Title:       "Шаг 2: Осязание",
			Description: touchExample,
			Emoji:       "🤚",
		},
		{
			Number:      groundingHearCount,
			Title:       "Шаг 3: Слух",
			Description: "Назовите 3 звука, которые вы слышите.\n\nНапример: шум улицы, тиканье часов, ваше дыхание.",
			Emoji:       "👂",
		},
		{
			Number:      groundingSmellCount,
			Title:       "Шаг 4: Обоняние",
			Description: "Назовите 2 запаха, которые вы чувствуете или любите.\n\nНапример: кофе, свежий воздух, аромат цветов.",
			Emoji:       "👃",
		},
		{
			Number:      groundingTasteCount,
			Title:       "Шаг 5: Вкус",
			Description: "Назовите 1 вкус, который вы ощущаете во рту.\n\nЕсли ничего нет — вспомните любимый вкус.",
			Emoji:       "👅",
		},
	}
}
