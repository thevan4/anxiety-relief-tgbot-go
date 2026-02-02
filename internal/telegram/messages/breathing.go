package messages

import "fmt"

const BreathingIntro = `🌬️ *Дыхание за 2 минуты*

Простое упражнение для успокоения нервной системы.

*Инструкция:*
1️⃣ Сядьте удобно, закройте глаза
2️⃣ Вдох через нос — 4 секунды
3️⃣ Задержка дыхания — 4 секунды
4️⃣ Выдох через рот — 6 секунд
5️⃣ Повторите 8 циклов (~2 минуты)

Я буду напоминать каждый этап. Готовы начать?`

func BreathingCycleText(cycleNum, totalCycles int, instruction string) string {
	return fmt.Sprintf("🌬️ *Цикл %d из %d*\n\n%s", cycleNum, totalCycles, instruction)
}

const BreathingCompletion = `✅ *Отлично!*

Вы завершили дыхательное упражнение.
Как вы себя чувствуете?`

const BreathingThanks = `✨ *Спасибо за практику!*

Регулярные упражнения помогают снизить уровень тревожности.

Выберите технику из меню ниже.`
