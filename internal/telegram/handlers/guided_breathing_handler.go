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
)

type GuidedBreathingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewGuidedBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *GuidedBreathingHandler {
	return &GuidedBreathingHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

func (h *GuidedBreathingHandler) Handle(ctx *th.Context, update telego.Update) error {
	if update.Message != nil {
		return h.handleMessage(ctx, update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}
	return nil
}

func (h *GuidedBreathingHandler) handleMessage(_ *th.Context, msg *telego.Message) error {
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
	h.showPatternSelection(h.ctx, msg.Chat.ID, userID)
	return nil
}

func (h *GuidedBreathingHandler) handleCallback(ctx *th.Context, cb *telego.CallbackQuery) error {
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

func (h *GuidedBreathingHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case strings.HasPrefix(data, "gbreath_pattern_"):
		patternIdx := strings.TrimPrefix(data, "gbreath_pattern_")
		idx, err := strconv.Atoi(patternIdx)
		if err == nil {
			h.startPattern(h.ctx, chatID, userID, messageID, idx)
		}
	case data == "gbreath_stop":
		h.stopBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_complete":
		h.completeBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_cancel":
		h.cancelBreathing(h.ctx, chatID, userID, messageID)
	}
}

func (h *GuidedBreathingHandler) showPatternSelection(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	patterns := techniques.GetBreathingPatterns()
	text := "🧘 *Управляемое дыхание*\n\nВыберите технику дыхания:\n\n"
	for _, p := range patterns {
		text += fmt.Sprintf("%s *%s*\n_%s_\n\n", p.Emoji, p.Name, p.Description)
	}

	buttons := h.buildPatternButtons(patterns)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	removeKeyboard := &telego.ReplyKeyboardRemove{RemoveKeyboard: true}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		text,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard)); err != nil {
		log.Printf("ERROR: send guided breathing intro: %v", err)
		return
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"_Используйте кнопки выше_",
	).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard)); err != nil {
		log.Printf("ERROR: send guided breathing use buttons: %v", err)
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
		{Text: "❌ Отмена", CallbackData: "gbreath_cancel"},
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

	text := fmt.Sprintf("%s *%s*\n\n*Цикл %d из %d*\n\n%s %s\n\n_%d сек_",
		pattern.Emoji, pattern.Name, cycle, pattern.Cycles,
		phase.Emoji, phase.Name, int(phase.Duration.Seconds()))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: "🛑 Стоп", CallbackData: "gbreath_stop"}},
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

	text := fmt.Sprintf(`✅ *Отлично!*

Вы завершили упражнение "%s".

Как вы себя чувствуете?`, pattern.Name)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "✅ Лучше", CallbackData: "gbreath_complete"},
				{Text: "🔄 Повторить", CallbackData: fmt.Sprintf("gbreath_pattern_%d", patternIdx)},
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

func (h *GuidedBreathingHandler) stopBreathing(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	patterns := techniques.GetBreathingPatterns()
	text := "🧘 *Управляемое дыхание*\n\nУпражнение остановлено. Выберите технику:\n\n"
	for _, p := range patterns {
		text += fmt.Sprintf("%s *%s*\n_%s_\n\n", p.Emoji, p.Name, p.Description)
	}

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

func (h *GuidedBreathingHandler) completeBreathing(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	text := `✨ *Спасибо за практику!*

Регулярные дыхательные упражнения помогают снизить уровень тревожности и улучшить концентрацию.`

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}); err != nil {
		log.Printf("ERROR: edit guided breathing complete: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send guided breathing menu: %v", err)
	}
}

func (h *GuidedBreathingHandler) cancelBreathing(
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
		log.Printf("ERROR: edit guided breathing cancel: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send guided breathing menu: %v", err)
	}
}
