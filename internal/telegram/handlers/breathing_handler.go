package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
)

type BreathingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *BreathingHandler {
	return &BreathingHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

func (h *BreathingHandler) Handle(ctx *th.Context, update telego.Update) error {
	if update.Message != nil {
		return h.handleBreathingMessage(update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleBreathingCallback(ctx, update.CallbackQuery)
	}
	return nil
}

func (h *BreathingHandler) handleBreathingMessage(msg *telego.Message) error {
	userID := msg.From.ID
	h.rateLimiter.WaitAndGo(h.ctx, userID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}
	h.statistics.IncreaseRequestsStatisticForUser(
		userID,
		msg.From.Username,
		msg.From.IsPremium,
		msg.From.IsBot,
	)
	h.startBreathingExercise(h.ctx, msg.Chat.ID, userID)
	return nil
}

func (h *BreathingHandler) handleBreathingCallback(ctx *th.Context, cb *telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		log.Printf("ERROR: callback query message is inaccessible")
		return nil
	}
	chatID := msg.Chat.ID
	userID := cb.From.ID
	messageID := msg.MessageID
	h.rateLimiter.WaitAndGo(h.ctx, userID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}

	h.processCallback(chatID, userID, messageID, cb.Data)

	if err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	}); err != nil {
		log.Printf("ERROR: answer callback query: %v", err)
		return err
	}
	return nil
}

func (h *BreathingHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch data {
	case "breathing_complete":
		h.completeBreathing(h.ctx, chatID, userID, messageID)
	case "breathing_cancel":
		h.cancelBreathing(h.ctx, chatID, userID, messageID)
	case "breathing_stop":
		h.stopBreathing(h.ctx, chatID, userID, messageID)
	case "breathing_start":
		if err := h.sessionStorage.SetState(h.ctx, userID, session.StateBreathingRunning); err != nil {
			log.Printf("ERROR: set state breathing running: %v", err)
		}
		h.runBreathingCycle(h.ctx, chatID, userID, messageID)
	}
}

func (h *BreathingHandler) startBreathingExercise(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateBreathingActive); err != nil {
		log.Printf("ERROR: set state breathing active: %v", err)
	}
	h.showBreathingIntro(ctx, chatID, 0)
}

func (h *BreathingHandler) showBreathingIntro(ctx context.Context, chatID int64, messageID int) {
	intro := `🌬️ *Дыхание за 2 минуты*

Простое упражнение для успокоения нервной системы.

*Инструкция:*
1️⃣ Сядьте удобно, закройте глаза
2️⃣ Вдох через нос — 4 секунды
3️⃣ Задержка дыхания — 4 секунды
4️⃣ Выдох через рот — 6 секунд
5️⃣ Повторите 8 циклов (~2 минуты)

Я буду напоминать каждый этап. Готовы начать?`

	inlineKeyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "breathing_start"},
				{Text: "❌ Отмена", CallbackData: "breathing_cancel"},
			},
		},
	}

	if messageID == 0 {
		removeKeyboard := &telego.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		}

		_, err := h.bot.SendMessage(ctx, tu.Message(
			tu.ID(chatID),
			intro,
		).WithParseMode("Markdown").WithReplyMarkup(inlineKeyboard))
		if err != nil {
			log.Printf("ERROR: send breathing intro: %v", err)
			return
		}

		if _, err := h.bot.SendMessage(ctx, tu.Message(
			tu.ID(chatID),
			"_Используйте кнопки выше_",
		).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard)); err != nil {
			log.Printf("ERROR: send breathing use buttons: %v", err)
		}
	} else {
		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        intro,
			ParseMode:   "Markdown",
			ReplyMarkup: inlineKeyboard,
		}); err != nil {
			log.Printf("ERROR: edit breathing intro: %v", err)
		}
	}
}

const breathingStepsPerCycle = 3
const breathingTotalCycles = 8

func (h *BreathingHandler) runBreathingCycle(
	ctx context.Context, chatID int64, userID int64, messageID int,
) {
	cycles := techniques.GetBreathingCycles()
	for i, cycle := range cycles {
		if !h.sendCycleStep(ctx, chatID, userID, messageID, i, cycle) {
			return
		}
		timer := time.NewTimer(cycle.Duration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
	h.sendBreathingCompletion(ctx, chatID, userID, messageID)
}

func (h *BreathingHandler) sendCycleStep(
	ctx context.Context, chatID, userID int64, messageID int, stepIndex int, cycle techniques.BreathingCycle,
) bool {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateBreathingRunning {
		return false
	}
	cycleNum := stepIndex/breathingStepsPerCycle + 1
	text := fmt.Sprintf("🌬️ *Цикл %d из %d*\n\n%s", cycleNum, breathingTotalCycles, cycle.Instruction)
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: "🛑 Стоп", CallbackData: "breathing_stop"}},
		},
	}
	_, err = h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Printf("ERROR: edit breathing message: %v", err)
	}
	return true
}

func (h *BreathingHandler) sendBreathingCompletion(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateBreathingRunning {
		return
	}
	completionText := `✅ *Отлично!*

Вы завершили дыхательное упражнение.
Как вы себя чувствуете?`
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "✅ Лучше", CallbackData: "breathing_complete"},
				{Text: "🔄 Повторить", CallbackData: "breathing_start"},
			},
		},
	}
	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        completionText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: send completion message: %v", err)
	}
}

func (h *BreathingHandler) stopBreathing(ctx context.Context, chatID int64, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateBreathingActive); err != nil {
		log.Printf("ERROR: set state breathing active: %v", err)
	}
	h.showBreathingIntro(ctx, chatID, messageID)
}

func (h *BreathingHandler) completeBreathing(
	ctx context.Context, chatID int64, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	text := `✨ *Спасибо за практику!*

Регулярные упражнения помогают снизить уровень тревожности.`

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}); err != nil {
		log.Printf("ERROR: edit breathing complete: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send breathing menu: %v", err)
	}
}

func (h *BreathingHandler) cancelBreathing(
	ctx context.Context, chatID int64, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      "❌ Упражнение отменено.",
	}); err != nil {
		log.Printf("ERROR: edit breathing cancel: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send breathing menu: %v", err)
	}
}
