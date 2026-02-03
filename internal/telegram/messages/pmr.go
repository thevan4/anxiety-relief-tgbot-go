package messages

import "fmt"

func PMRIntro(muscleGroupsCount int) string {
	return fmt.Sprintf(`💪 *Прогрессивная мышечная релаксация*

Техника глубокого расслабления через напряжение и расслабление мышц.

*Как это работает:*
1. Напрягаете группу мышц на 7 секунд
2. Расслабляете на 15 секунд
3. Переходите к следующей группе

*%d групп мышц* — от рук до ног.

Готовы начать?`, muscleGroupsCount)
}

const pmrProgressBarWidth = 10

// PMRTensePhaseWithProgress returns formatted text for tense phase with progress bar.
func PMRTensePhaseWithProgress(instruction string, elapsed, totalSec int) string {
	progress := ProgressBar(elapsed, totalSec, pmrProgressBarWidth)
	return fmt.Sprintf("🔴 *НАПРЯГИТЕ*\n\n%s\n\n%s", instruction, progress)
}

// PMRRelaxPhaseWithProgress returns formatted text for relax phase with progress bar.
func PMRRelaxPhaseWithProgress(instruction string, elapsed, totalSec int) string {
	progress := ProgressBar(elapsed, totalSec, pmrProgressBarWidth)
	return fmt.Sprintf("🟢 *РАССЛАБЬТЕ*\n\n%s\n\n%s", instruction, progress)
}

// Legacy functions without progress bar.
func PMRTensePhase(instruction string) string {
	return fmt.Sprintf("🔴 *НАПРЯГИТЕ*\n\n%s\n\n_7 секунд..._", instruction)
}

func PMRRelaxPhase(instruction string) string {
	return fmt.Sprintf("🟢 *РАССЛАБЬТЕ*\n\n%s\n\n_15 секунд..._", instruction)
}

func PMRMuscleStep(emoji, name string, number, total int, phaseText string) string {
	return fmt.Sprintf("%s *%s* (%d/%d)\n\n%s", emoji, name, number, total, phaseText)
}

const PMRCompletion = `✅ *Отлично!*

Вы завершили прогрессивную мышечную релаксацию.

Ваше тело теперь полностью расслаблено. Посидите ещё минуту, наслаждаясь этим состоянием.`

const PMRStopped = `💪 *Прогрессивная мышечная релаксация*

Упражнение остановлено.

Хотите начать сначала?`

const PMRThanks = `✨ *Спасибо за практику!*

Прогрессивная мышечная релаксация снижает мышечное напряжение и уровень стресса.

Выберите технику из меню ниже.`
