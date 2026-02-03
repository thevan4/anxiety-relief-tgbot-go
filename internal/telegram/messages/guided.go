package messages

import "fmt"

func GuidedBreathingIntro(patternsInfo string) string {
	return "🧘 *Управляемое дыхание*\n\nВыберите технику дыхания:\n\n" + patternsInfo
}

func GuidedBreathingStoppedIntro(patternsInfo string) string {
	return "🧘 *Управляемое дыхание*\n\nУпражнение остановлено. Выберите технику:\n\n" + patternsInfo
}

const guidedProgressBarWidth = 10

// GuidedBreathingPhase returns formatted text for guided breathing phase with progress bar.
func GuidedBreathingPhase(patternEmoji, patternName string, cycle, totalCycles int, phaseEmoji, phaseName string, elapsed, totalSec int) string {
	progress := ProgressBar(elapsed, totalSec, guidedProgressBarWidth)
	return fmt.Sprintf("%s *%s*\n\n*Цикл %d/%d*\n%s %s\n%s",
		patternEmoji, patternName, cycle, totalCycles,
		phaseName, phaseEmoji, progress)
}

func GuidedBreathingCompletion(patternName string) string {
	return fmt.Sprintf(`✅ *Отлично!*

Вы завершили упражнение "%s".

Как вы себя чувствуете?`, patternName)
}

const GuidedBreathingThanks = `✨ *Спасибо за практику!*

Регулярные дыхательные упражнения помогают снизить уровень тревожности и улучшить концентрацию.

Выберите технику из меню ниже.`
