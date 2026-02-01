package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
)

type GroundingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewGroundingHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *GroundingHandler {
	return &GroundingHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

func (h *GroundingHandler) Handle(ctx *th.Context, update telego.Update) error {
	if update.Message != nil {
		return h.handleGroundingMessage(update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleGroundingCallback(ctx, update.CallbackQuery)
	}
	return nil
}

func (h *GroundingHandler) handleGroundingMessage(msg *telego.Message) error {
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
	h.startGroundingExercise(h.ctx, msg.Chat.ID, userID)
	return nil
}

func (h *GroundingHandler) handleGroundingCallback(ctx *th.Context, cb *telego.CallbackQuery) error {
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
	h.applyGroundingCallbackAction(h.ctx, chatID, userID, messageID, cb.Data)
	if err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	}); err != nil {
		log.Printf("ERROR: answer callback query: %v", err)
		return err
	}
	return nil
}

func (h *GroundingHandler) applyGroundingCallbackAction(
	ctx context.Context, chatID, userID int64, messageID int, data string,
) {
	switch {
	case strings.HasPrefix(data, "grounding_step_"):
		stepNum := strings.TrimPrefix(data, "grounding_step_")
		if stepNum == "1" || stepNum == "2" || stepNum == "3" || stepNum == "4" || stepNum == "5" {
			h.showGroundingStep(ctx, chatID, userID, messageID, stepNum)
		}
	case data == "grounding_complete":
		h.completeGrounding(ctx, chatID, userID, messageID)
	case data == "grounding_cancel":
		h.cancelGrounding(ctx, chatID, userID, messageID)
	}
}

func (h *GroundingHandler) startGroundingExercise(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGroundingStep1); err != nil {
		log.Printf("ERROR: set state grounding step1: %v", err)
	}

	intro := `🌿 *Якорение 5-4-3-2-1*

Техника для возвращения в настоящий момент.

*Суть:* Назовите вслух или про себя:

• 5 вещей, которые видите 👁️
• 4 вещи, которые ощущаете 🤚
• 3 звука, которые слышите 👂
• 2 запаха, которые чувствуете 👃
• 1 вкус во рту 👅

Готовы начать?`

	inlineKeyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "grounding_step_1"},
				{Text: "❌ Отмена", CallbackData: "grounding_cancel"},
			},
		},
	}

	removeKeyboard := &telego.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}

	_, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		intro,
	).WithParseMode("Markdown").WithReplyMarkup(inlineKeyboard))
	if err != nil {
		log.Printf("ERROR: send grounding intro: %v", err)
		return
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"_Используйте кнопки выше_",
	).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard)); err != nil {
		log.Printf("ERROR: send grounding use buttons: %v", err)
	}
}

func getGroundingStepConfig() []struct {
	state    session.State
	nextStep string
} {
	return []struct {
		state    session.State
		nextStep string
	}{
		{session.StateGroundingStep1, "2"},
		{session.StateGroundingStep2, "3"},
		{session.StateGroundingStep3, "4"},
		{session.StateGroundingStep4, "5"},
		{session.StateGroundingStep5, "complete"},
	}
}

func (h *GroundingHandler) showGroundingStep(
	ctx context.Context, chatID, userID int64, messageID int, stepNum string,
) {
	if len(stepNum) != 1 {
		return
	}
	cfg := getGroundingStepConfig()
	idx := int(stepNum[0] - '1')
	if idx < 0 || idx >= len(cfg) {
		return
	}
	stepCfg := cfg[idx]
	if err := h.sessionStorage.SetState(ctx, userID, stepCfg.state); err != nil {
		log.Printf("ERROR: set state grounding: %v", err)
	}

	steps := techniques.GetGroundingSteps()
	step := steps[idx]
	nextStep := stepCfg.nextStep

	text := fmt.Sprintf("%s *%s*\n\n%s", step.Emoji, step.Title, step.Description)

	var keyboard *telego.InlineKeyboardMarkup
	if nextStep == "complete" {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: "✅ Завершить", CallbackData: "grounding_complete"},
					{Text: "❌ Отмена", CallbackData: "grounding_cancel"},
				},
			},
		}
	} else {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: "➡️ Далее", CallbackData: "grounding_step_" + nextStep},
					{Text: "❌ Отмена", CallbackData: "grounding_cancel"},
				},
			},
		}
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit grounding message: %v", err)
	}
}

func (h *GroundingHandler) completeGrounding(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	text := `✨ *Отлично!*

Вы завершили технику якорения 5-4-3-2-1.

Эта практика помогает выйти из тревожных мыслей и вернуться в настоящий момент.`

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}); err != nil {
		log.Printf("ERROR: edit grounding complete: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send grounding menu: %v", err)
	}
}

func (h *GroundingHandler) cancelGrounding(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      "❌ Упражнение отменено.",
	}); err != nil {
		log.Printf("ERROR: edit grounding cancel: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send grounding menu: %v", err)
	}
}
