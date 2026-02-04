package techniques

// ThoughtCategory represents a category of anxious thoughts for labeling.
type ThoughtCategory struct {
	ID          string
	Name        string
	Emoji       string
	Description string
	Example     string
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
			ID:          d.id,
			Name:        d.name,
			Emoji:       d.emoji,
			Description: d.desc,
			Example:     d.example,
		}
	}
	return result
}

type categoryData struct {
	id      string
	name    string
	emoji   string
	desc    string
	example string
}

func getThoughtCategoryData() []categoryData {
	return []categoryData{
		{"worry", "Беспокойство", "😰",
			"Тревога о будущих событиях или исходах",
			"«А что если я провалю собеседование?»"},
		{"catastrophic", "Катастрофизация", "🌪️",
			"Представление наихудшего сценария",
			"«Если я ошибусь, всё будет ужасно!»"},
		{"self_doubt", "Самосомнение", "🤔",
			"Сомнения в своих способностях",
			"«Я недостаточно хорош для этого»"},
		{"perfectionist", "Перфекционизм", "🎯",
			"Нереалистично высокие стандарты",
			"«Если не идеально — значит плохо»"},
		{"comparison", "Сравнение", "⚖️",
			"Сравнение себя с другими",
			"«Все справляются лучше меня»"},
		{"rumination", "Руминация", "🔄",
			"Постоянное возвращение к прошлым событиям",
			"«Почему я тогда так сказал?»"},
		{"control", "Контроль", "🎮",
			"Желание контролировать неконтролируемое",
			"«Я должен всё предусмотреть»"},
		{"rejection", "Страх отвержения", "💔",
			"Боязнь быть отвергнутым",
			"«Они подумают, что я странный»"},
		{"health", "Тревога о здоровье", "🏥",
			"Чрезмерное беспокойство о здоровье",
			"«Этот симптом — признак болезни?»"},
		{"social", "Социальная тревога", "👥",
			"Страх социальных ситуаций",
			"«Все будут на меня смотреть»"},
		{"financial", "Финансовые переживания", "💰",
			"Стресс из-за денег",
			"«Что если денег не хватит?»"},
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
