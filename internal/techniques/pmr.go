package techniques

import "time"

// PMR timing constants.
const (
	tenseDuration = 7 * time.Second
	relaxDuration = 15 * time.Second
)

// MuscleGroup represents one muscle group for Progressive Muscle Relaxation.
type MuscleGroup struct {
	Number        int
	ID            string
	Emoji         string
	TenseDuration time.Duration
	RelaxDuration time.Duration
}

// GetMuscleGroups returns the muscle groups for Progressive Muscle Relaxation exercise.
func GetMuscleGroups() []MuscleGroup {
	return buildMuscleGroups()
}

func buildMuscleGroups() []MuscleGroup {
	groups := getMuscleGroupData()
	result := make([]MuscleGroup, len(groups))
	for i, g := range groups {
		result[i] = MuscleGroup{
			Number:        i + 1,
			ID:            g.id,
			Emoji:         g.emoji,
			TenseDuration: tenseDuration,
			RelaxDuration: relaxDuration,
		}
	}
	return result
}

type muscleData struct {
	id    string
	emoji string
}

func getMuscleGroupData() []muscleData {
	return append(getUpperBodyMuscles(), getLowerBodyMuscles()...)
}

func getUpperBodyMuscles() []muscleData {
	return []muscleData{
		{id: "hands", emoji: "✊"},
		{id: "forearms", emoji: "💪"},
		{id: "forehead", emoji: "😤"},
		{id: "eyes", emoji: "😣"},
		{id: "jaw", emoji: "😬"},
	}
}

func getLowerBodyMuscles() []muscleData {
	return []muscleData{
		{id: "neck", emoji: "🙆"},
		{id: "chest", emoji: "🫁"},
		{id: "stomach", emoji: "🎯"},
		{id: "thighs", emoji: "🦵"},
		{id: "calves", emoji: "🦶"},
	}
}
