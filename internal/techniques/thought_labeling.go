package techniques

// ThoughtCategory represents a category of anxious thoughts for labeling.
type ThoughtCategory struct {
	ID    string
	Emoji string
}

// GetThoughtCategories returns the list of thought categories for labeling exercise.
func GetThoughtCategories() []ThoughtCategory {
	return buildThoughtCategories()
}

func buildThoughtCategories() []ThoughtCategory {
	data := getThoughtCategoryData()
	result := make([]ThoughtCategory, len(data))
	for i, d := range data {
		result[i] = ThoughtCategory{
			ID:    d.id,
			Emoji: d.emoji,
		}
	}
	return result
}

type categoryData struct {
	id    string
	emoji string
}

func getThoughtCategoryData() []categoryData {
	return []categoryData{
		{"worry", "😰"},
		{"catastrophic", "🌪️"},
		{"self_doubt", "🤔"},
		{"perfectionist", "🎯"},
		{"comparison", "⚖️"},
		{"rumination", "🔄"},
		{"control", "🎮"},
		{"rejection", "💔"},
		{"health", "🏥"},
		{"social", "👥"},
		{"financial", "💰"},
	}
}

// LabelingSession tracks user's progress in thought labeling exercise.
type LabelingSession struct {
	ThoughtsLabeled int
	Categories      map[string]int
}

// NewLabelingSession creates a new thought labeling session tracker.
func NewLabelingSession() *LabelingSession {
	return &LabelingSession{
		ThoughtsLabeled: 0,
		Categories:      make(map[string]int),
	}
}

// AddLabel records a labeled thought in the session.
func (s *LabelingSession) AddLabel(categoryID string) {
	s.ThoughtsLabeled++
	s.Categories[categoryID]++
}

// GetMostFrequent returns the most frequently labeled thought category.
func (s *LabelingSession) GetMostFrequent() string {
	maxCount := 0
	mostFrequent := ""
	for id, count := range s.Categories {
		if count > maxCount {
			maxCount = count
			mostFrequent = id
		}
	}
	return mostFrequent
}
