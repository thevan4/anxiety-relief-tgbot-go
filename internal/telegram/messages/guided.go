package messages

import "fmt"

func GuidedBreathingIntro(patternsInfo string) string {
	return "🧘 *Управляемое дыхание*\n\nВыберите технику дыхания:\n\n" + patternsInfo
}

func GuidedBreathingStoppedIntro(patternsInfo string) string {
	return "🧘 *Управляемое дыхание*\n\nУпражнение остановлено. Выберите технику:\n\n" + patternsInfo
}

func GuidedBreathingPhase(patternEmoji, patternName string, cycle, totalCycles int, phaseEmoji, phaseName string, durationSec int) string {
	return fmt.Sprintf("%s *%s*\n\n*Цикл %d из %d*\n\n%s %s\n\n_%d сек_",
		patternEmoji, patternName, cycle, totalCycles,
		phaseEmoji, phaseName, durationSec)
}

func GuidedBreathingCompletion(patternName string) string {
	return fmt.Sprintf(`✅ *Отлично!*

Вы завершили упражнение "%s".

Как вы себя чувствуете?`, patternName)
}

const GuidedBreathingThanks = `✨ *Спасибо за практику!*

Регулярные дыхательные упражнения помогают снизить уровень тревожности и улучшить концентрацию.

Выберите технику из меню ниже.`
