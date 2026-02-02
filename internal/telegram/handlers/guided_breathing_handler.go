package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
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

type GuidedBreathingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
	sessionManager *session.SessionManager
}

func NewGuidedBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
) *GuidedBreathingHandler {
	return &GuidedBreathingHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
		sessionManager: sessionManager,
	}
}

// HandleMenuSelect handles selection from main menu
func (h *GuidedBreathingHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		return nil
	}
	userID := cb.From.ID
	chatID := msg.Chat.ID
	messageID := msg.MessageID

	// Answer callback FIRST; if too old - delete message and stop
	if !AnswerCallbackOrDelete(h.ctx, h.bot, cb.ID, chatID, messageID) {
		return nil
	}

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

	h.showPatternSelection(h.ctx, chatID, userID, messageID)
	return nil
}

// HandleCallback handles guided breathing callbacks
func (h *GuidedBreathingHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		log.Printf("ERROR: callback query message is inaccessible")
		return nil
	}
	chatID := msg.Chat.ID
	userID := cb.From.ID
	messageID := msg.MessageID

	// Answer callback FIRST; if too old - delete message and stop
	if !AnswerCallbackOrDelete(h.ctx, h.bot, cb.ID, chatID, messageID) {
		return nil
	}

	h.rateLimiter.WaitAndGo(h.ctx, userID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}

	h.processCallback(chatID, userID, messageID, cb.Data)
	return nil
}

func (h *GuidedBreathingHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case strings.HasPrefix(data, "gbreath_pattern_"):
		patternIdx := strings.TrimPrefix(data, "gbreath_pattern_")
		idx, err := strconv.Atoi(patternIdx)
		if err == nil {
			// Create new session context (cancels previous if any)
			sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
			go h.startPattern(sessionCtx, chatID, userID, messageID, idx)
		}
	case data == "gbreath_stop":
		h.sessionManager.CancelSession(userID)
		h.stopBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_complete":
		h.sessionManager.CancelSession(userID)
		h.completeBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_cancel":
		h.sessionManager.CancelSession(userID)
		h.cancelBreathing(h.ctx, chatID, userID, messageID)
	}
}

func (h *GuidedBreathingHandler) buildPatternsInfo(patterns []techniques.BreathingPattern) string {
	var info string
	for _, p := range patterns {
		info += fmt.Sprintf("%s *%s*\n_%s_\n\n", p.Emoji, p.Name, p.Description)
	}
	return info
}

func (h *GuidedBreathingHandler) showPatternSelection(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	patterns := techniques.GetBreathingPatterns()
	text := messages.GuidedBreathingIntro(h.buildPatternsInfo(patterns))

	buttons := h.buildPatternButtons(patterns)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit guided breathing intro: %v", err)
	}
}

func (h *GuidedBreathingHandler) buildPatternButtons(
	patterns []techniques.BreathingPattern,
) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	for i, p := range patterns {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			{Text: fmt.Sprintf("%s %s", p.Emoji, p.Name), CallbackData: fmt.Sprintf("gbreath_pattern_%d", i)},
		})
	}
	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: messages.Cancel, CallbackData: "gbreath_cancel"},
	})
	return buttons
}

func (h *GuidedBreathingHandler) startPattern(
	ctx context.Context, chatID, userID int64, messageID int, patternIdx int,
) {
	patterns := techniques.GetBreathingPatterns()
	if patternIdx < 0 || patternIdx >= len(patterns) {
		return
	}
	pattern := patterns[patternIdx]

	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingRunning); err != nil {
		log.Printf("ERROR: set state guided breathing running: %v", err)
	}

	if !h.runBreathingCycles(ctx, chatID, userID, messageID, pattern) {
		return
	}

	h.sendCompletion(ctx, chatID, userID, messageID, patternIdx)
}

func (h *GuidedBreathingHandler) runBreathingCycles(
	ctx context.Context, chatID, userID int64, messageID int, pattern techniques.BreathingPattern,
) bool {
	phases := techniques.GetBreathingPhases(pattern)

	for cycle := 1; cycle <= pattern.Cycles; cycle++ {
		for _, phase := range phases {
			if !h.runPhase(ctx, chatID, userID, messageID, pattern, cycle, phase) {
				return false
			}
		}
	}
	return true
}

func (h *GuidedBreathingHandler) runPhase(
	ctx context.Context, chatID, userID int64, messageID int,
	pattern techniques.BreathingPattern, cycle int, phase techniques.BreathingPhase,
) bool {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateGuidedBreathingRunning {
		return false
	}

	text := messages.GuidedBreathingPhase(
		pattern.Emoji, pattern.Name, cycle, pattern.Cycles,
		phase.Emoji, phase.Name, int(phase.Duration.Seconds()),
	)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: messages.Stop, CallbackData: "gbreath_stop"}},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit guided breathing message: %v", err)
	}

	timer := time.NewTimer(phase.Duration)
	select {
	case <-ctx.Done():
		timer.Stop()
		return false
	case <-timer.C:
		return true
	}
}

func (h *GuidedBreathingHandler) sendCompletion(
	ctx context.Context, chatID, userID int64, messageID int, patternIdx int,
) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateGuidedBreathingRunning {
		return
	}

	patterns := techniques.GetBreathingPatterns()
	pattern := patterns[patternIdx]

	text := messages.GuidedBreathingCompletion(pattern.Name)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.FeelBetter, CallbackData: "gbreath_complete"},
				{Text: messages.Repeat, CallbackData: fmt.Sprintf("gbreath_pattern_%d", patternIdx)},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: send guided breathing completion: %v", err)
	}
}

func (h *GuidedBreathingHandler) stopBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	patterns := techniques.GetBreathingPatterns()
	text := messages.GuidedBreathingStoppedIntro(h.buildPatternsInfo(patterns))

	buttons := h.buildPatternButtons(patterns)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit guided breathing stop: %v", err)
	}
}

func (h *GuidedBreathingHandler) completeBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.GuidedBreathingThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit guided breathing complete: %v", err)
	}
}

func (h *GuidedBreathingHandler) cancelBreathing(ctx context.Context, chatID, userID int64, messageID int) {
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
		if !HandleEditError(ctx, h.bot, err, chatID, messageID) {
			log.Printf("ERROR: edit guided breathing cancel: %v", err)
		}
	}
}
