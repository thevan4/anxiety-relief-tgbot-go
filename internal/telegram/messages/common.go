package messages

import (
	"fmt"
	"strings"
)

const (
	FeelBetter = "✅ Лучше"
	Repeat     = "🔄 Повторить"
	Start      = "▶️ Начать"
	Cancel     = "❌ Отмена"
	Stop       = "🛑 Стоп"
	Next       = "➡️ Далее"
	Complete   = "✅ Завершить"
	Done       = "✅ Готово"
	BackToMenu = "🏠 В меню"

	MainMenuText = `Выберите технику для работы с тревожностью:

*Быстрые (2-5 мин):*
🌬️ Дыхание — успокоение нервной системы
🌿 Якорение — возврат в настоящий момент

*Продвинутые (5-15 мин):*
🧘 Управляемое дыхание — разные паттерны
💪 Мышечная релаксация — снятие напряжения
🏷️ Маркировка мыслей — работа с тревогой
🌅 Визуализация — расслабление`

	progressFilled = "▓"
	progressEmpty  = "░"
)

// ProgressBar generates a progress bar string.
// current - current progress value
// total - maximum value
// width - number of characters in the bar
func ProgressBar(current, total, width int) string {
	if total <= 0 || width <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

	filled := (current * width) / total
	empty := width - filled

	return strings.Repeat(progressFilled, filled) + strings.Repeat(progressEmpty, empty)
}

// TimerCountdown returns emoji timer with remaining seconds.
// elapsed - seconds passed
// total - total seconds
// Returns format like "⏱️ 4" or "✨" when done
func TimerCountdown(elapsed, total int) string {
	remaining := total - elapsed
	if remaining <= 0 {
		return "✨"
	}
	return fmt.Sprintf("⏱️ %d", remaining)
}
