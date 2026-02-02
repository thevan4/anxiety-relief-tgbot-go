package handlers

import (
	"context"
	"log"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
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

// HandleMenuSelect handles selection from main menu
func (h *BreathingHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
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

	h.showBreathingIntro(h.ctx, chatID, userID, messageID)
	return nil
}

// HandleCallback handles breathing exercise callbacks
func (h *BreathingHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
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

func (h *BreathingHandler) showBreathingIntro(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateBreathingActive); err != nil {
		log.Printf("ERROR: set state breathing active: %v", err)
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.Start, CallbackData: "breathing_start"},
				{Text: messages.Cancel, CallbackData: "breathing_cancel"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.BreathingIntro,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit breathing intro: %v", err)
	}
}

const breathingStepsPerCycle = 3
const breathingTotalCycles = 8

func (h *BreathingHandler) runBreathingCycle(ctx context.Context, chatID, userID int64, messageID int) {
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
	text := messages.BreathingCycleText(cycleNum, breathingTotalCycles, cycle.Instruction)
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: messages.Stop, CallbackData: "breathing_stop"}},
		},
	}
	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit breathing message: %v", err)
	}
	return true
}

func (h *BreathingHandler) sendBreathingCompletion(ctx context.Context, chatID, userID int64, messageID int) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateBreathingRunning {
		return
	}
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.FeelBetter, CallbackData: "breathing_complete"},
				{Text: messages.Repeat, CallbackData: "breathing_start"},
			},
		},
	}
	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.BreathingCompletion,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: send completion message: %v", err)
	}
}

func (h *BreathingHandler) stopBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateBreathingActive); err != nil {
		log.Printf("ERROR: set state breathing active: %v", err)
	}
	h.showBreathingIntro(ctx, chatID, userID, messageID)
}

func (h *BreathingHandler) completeBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	// Show thanks message then return to main menu
	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.BreathingThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit breathing complete: %v", err)
	}
}

func (h *BreathingHandler) cancelBreathing(ctx context.Context, chatID, userID int64, messageID int) {
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
		log.Printf("ERROR: edit breathing cancel: %v", err)
	}
}
