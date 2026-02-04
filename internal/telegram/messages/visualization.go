package messages

import "fmt"

func VisualizationIntro(scenesInfo string) string {
	return "🌅 *Мирная визуализация*\n\nТехника расслабления через воображение спокойных мест.\n\n*Выберите сцену:*\n\n" + scenesInfo
}

func VisualizationStoppedIntro(scenesInfo string) string {
	return "🌅 *Мирная визуализация*\n\nУпражнение остановлено. Выберите другую сцену:\n\n" + scenesInfo
}

func VisualizationSceneIntro(emoji, name, description, atmosphere string) string {
	return fmt.Sprintf(`%s *%s*

_%s_

🌿 *Атмосфера:* %s

Закройте глаза и погрузитесь в эту сцену...`, emoji, name, description, atmosphere)
}

// VisualizationStep returns formatted text for visualization step (legacy).
func VisualizationStep(emoji, name string, stepNum, totalSteps int, instruction string) string {
	return fmt.Sprintf(`%s *%s*

*Шаг %d из %d*

%s`, emoji, name, stepNum, totalSteps, instruction)
}

// VisualizationStepWithProgress returns formatted text for visualization step with timer countdown.
func VisualizationStepWithProgress(emoji, name string, stepNum, totalSteps int, instruction string, elapsed, totalSec int) string {
	timer := TimerCountdown(elapsed, totalSec)
	return fmt.Sprintf(`%s *%s*

*Шаг %d/%d*

%s

%s`, emoji, name, stepNum, totalSteps, instruction, timer)
}

// VisualizationIntroWithProgress returns formatted text for scene intro with timer countdown.
func VisualizationIntroWithProgress(emoji, name, description, atmosphere string, elapsed, totalSec int) string {
	timer := TimerCountdown(elapsed, totalSec)
	return fmt.Sprintf(`%s *%s*

_%s_

🌿 *Атмосфера:* %s

Закройте глаза и погрузитесь в эту сцену...

%s`, emoji, name, description, atmosphere, timer)
}

func VisualizationCompletion(sceneName string) string {
	return fmt.Sprintf(`✅ *Отлично!*

Вы завершили визуализацию "%s".

Медленно возвращайтесь в реальность. Пошевелите пальцами, глубоко вдохните и откройте глаза.

Как вы себя чувствуете?`, sceneName)
}

const VisualizationThanks = `✨ *Спасибо за практику!*

Визуализация — мощная техника для снижения стресса и тревожности. Регулярная практика усиливает эффект.

Выберите технику из меню ниже.`
