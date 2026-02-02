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

func VisualizationStep(emoji, name string, stepNum, totalSteps int, instruction string) string {
	return fmt.Sprintf(`%s *%s*

*Шаг %d из %d*

%s`, emoji, name, stepNum, totalSteps, instruction)
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
