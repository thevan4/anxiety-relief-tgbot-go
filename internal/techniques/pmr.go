package techniques

import "time"

// PMR timing constants.
const (
	tenseDuration = 7 * time.Second
	relaxDuration = 15 * time.Second
)

// MuscleGroup represents one muscle group for Progressive Muscle Relaxation.
type MuscleGroup struct {
	Number           int
	Name             string
	Emoji            string
	TenseInstruction string
	RelaxInstruction string
	TenseDuration    time.Duration
	RelaxDuration    time.Duration
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
			Number:           i + 1,
			Name:             g.name,
			Emoji:            g.emoji,
			TenseInstruction: g.tense,
			RelaxInstruction: g.relax,
			TenseDuration:    tenseDuration,
			RelaxDuration:    relaxDuration,
		}
	}
	return result
}

type muscleData struct {
	name  string
	emoji string
	tense string
	relax string
}

func getMuscleGroupData() []muscleData {
	return append(getUpperBodyMuscles(), getLowerBodyMuscles()...)
}

func getUpperBodyMuscles() []muscleData {
	return []muscleData{
		{
			name:  "Кисти рук",
			emoji: "✊",
			tense: "Сожмите кулаки как можно сильнее. Почувствуйте напряжение в пальцах и ладонях.",
			relax: "Разожмите кулаки и расслабьте руки. Почувствуйте тепло и расслабление.",
		},
		{
			name:  "Предплечья и бицепсы",
			emoji: "💪",
			tense: "Согните руки в локтях и напрягите бицепсы. Держите напряжение.",
			relax: "Опустите руки и полностью расслабьте их. Руки становятся тяжёлыми.",
		},
		{
			name:  "Лоб",
			emoji: "😤",
			tense: "Поднимите брови как можно выше, наморщите лоб. Почувствуйте напряжение.",
			relax: "Расслабьте лоб. Позвольте бровям опуститься. Лоб становится гладким.",
		},
		{
			name:  "Глаза и нос",
			emoji: "😣",
			tense: "Крепко зажмурьтесь и наморщите нос. Почувствуйте напряжение вокруг глаз.",
			relax: "Расслабьте глаза и нос. Веки становятся лёгкими и спокойными.",
		},
		{
			name:  "Челюсть",
			emoji: "😬",
			tense: "Сожмите челюсти и растяните губы в напряжённой улыбке.",
			relax: "Расслабьте челюсть, слегка приоткройте рот. Язык расслаблен.",
		},
	}
}

func getLowerBodyMuscles() []muscleData {
	return []muscleData{
		{
			name:  "Шея и плечи",
			emoji: "🙆",
			tense: "Поднимите плечи к ушам и напрягите шею. Держите напряжение.",
			relax: "Опустите плечи вниз и расслабьте шею. Почувствуйте облегчение.",
		},
		{
			name:  "Грудь и спина",
			emoji: "🫁",
			tense: "Глубоко вдохните, задержите дыхание и напрягите мышцы груди и спины.",
			relax: "Медленно выдохните и расслабьте грудь и спину. Дышите спокойно.",
		},
		{
			name:  "Живот",
			emoji: "🎯",
			tense: "Втяните живот и напрягите мышцы пресса. Держите напряжение.",
			relax: "Расслабьте живот. Позвольте ему свободно двигаться при дыхании.",
		},
		{
			name:  "Бёдра и ягодицы",
			emoji: "🦵",
			tense: "Напрягите ягодицы и бёдра, прижмите их к сиденью.",
			relax: "Расслабьте ягодицы и бёдра. Почувствуйте тяжесть в ногах.",
		},
		{
			name:  "Икры и стопы",
			emoji: "🦶",
			tense: "Потяните носки на себя, напрягите икры. Почувствуйте натяжение.",
			relax: "Расслабьте стопы и икры. Ноги становятся тёплыми и тяжёлыми.",
		},
	}
}
