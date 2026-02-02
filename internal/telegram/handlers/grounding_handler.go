package handlers

import (
	"context"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
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

// HandleMenuSelect handles selection from main menu
func (h *GroundingHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		return nil
	}
	userID := cb.From.ID
	chatID := msg.Chat.ID
	messageID := msg.MessageID

	h.rateLimiter.WaitAndGo(h.ctx, userID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}

	h.statistics.IncreaseRequestsStatisticForUser(
		userID,
		cb.From.Username,
		cb.From.IsPremium,
		cb.From.IsBot,
	)

	if err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	}); err != nil {
		log.Printf("ERROR: answer callback query: %v", err)
	}

	h.showGroundingIntro(h.ctx, chatID, userID, messageID)
	return nil
}

// HandleCallback handles grounding exercise callbacks
func (h *GroundingHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
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

func (h *GroundingHandler) showGroundingIntro(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGroundingStep1); err != nil {
		log.Printf("ERROR: set state grounding step1: %v", err)
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.Start, CallbackData: "grounding_step_1"},
				{Text: messages.Cancel, CallbackData: "grounding_cancel"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.GroundingIntro,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit grounding intro: %v", err)
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

	text := messages.GroundingStepText(step.Emoji, step.Title, step.Description)

	var keyboard *telego.InlineKeyboardMarkup
	if nextStep == "complete" {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: messages.Complete, CallbackData: "grounding_complete"},
					{Text: messages.Cancel, CallbackData: "grounding_cancel"},
				},
			},
		}
	} else {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: messages.Next, CallbackData: "grounding_step_" + nextStep},
					{Text: messages.Cancel, CallbackData: "grounding_cancel"},
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

func (h *GroundingHandler) completeGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.GroundingCompletion,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit grounding complete: %v", err)
	}
}

func (h *GroundingHandler) cancelGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit grounding cancel: %v", err)
	}
}
