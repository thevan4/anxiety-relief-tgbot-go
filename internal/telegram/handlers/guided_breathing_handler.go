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
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
)

// GuidedBreathingHandler handles advanced breathing patterns with different timings.
type GuidedBreathingHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.SessionManager
	callbackProcessor *CallbackProcessor
}

// NewGuidedBreathingHandler creates a new guided breathing exercise handler.
func NewGuidedBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
	callbackProcessor *CallbackProcessor,
) *GuidedBreathingHandler {
	return &GuidedBreathingHandler{
		ctx:               ctx,
		bot:               bot,
		localizer:         localizer,
		statistics:        statistics,
		sessionStorage:    sessionStorage,
		sessionManager:    sessionManager,
		callbackProcessor: callbackProcessor,
	}
}

// getLang returns user's language from session or default.
func (h *GuidedBreathingHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (h *GuidedBreathingHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.MenuBreathing, CallbackData: "menu_breathing"},
				{Text: m.MenuGrounding, CallbackData: "menu_grounding"},
			},
			{
				{Text: m.MenuGuided, CallbackData: "menu_guided"},
				{Text: m.MenuPMR, CallbackData: "menu_pmr"},
			},
			{
				{Text: m.MenuThought, CallbackData: "menu_thought"},
				{Text: m.MenuVisualization, CallbackData: "menu_visual"},
			},
			{
				{Text: m.MenuLang, CallbackData: "menu_lang"},
			},
		},
	}
}

// HandleMenuSelect handles selection from main menu.
func (h *GuidedBreathingHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showPatternSelection(h.ctx, info.ChatID, info.UserID, info.MessageID)
	return nil
}

// HandleCallback handles guided breathing callbacks.
func (h *GuidedBreathingHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.processCallback(info.ChatID, info.UserID, info.MessageID, info.Data)
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

// patternEmojis are fixed emojis for each breathing pattern.
var patternEmojis = []string{"📦", "😴", "⚡", "🚀"}

// getLocalizedPattern returns localized name and description for pattern by index.
func (h *GuidedBreathingHandler) getLocalizedPattern(idx int, m localization.Messages) (string, string) {
	switch idx {
	case 0:
		return m.PatternBoxName, m.PatternBoxDesc
	case 1:
		return m.PatternRelaxingName, m.PatternRelaxingDesc
	case 2:
		return m.PatternEnergizingName, m.PatternEnergizingDesc
	case 3:
		return m.PatternQuickName, m.PatternQuickDesc
	default:
		return "", ""
	}
}

// getLocalizedPhaseName returns localized phase name.
func (h *GuidedBreathingHandler) getLocalizedPhaseName(phaseName string, m localization.Messages) string {
	switch phaseName {
	case "Вдох через нос":
		return m.BreathingInhale
	case "Задержка":
		return m.BreathingHold
	case "Выдох через рот":
		return m.BreathingExhale
	case "Пауза":
		return m.BreathingPause
	default:
		return phaseName
	}
}

func (h *GuidedBreathingHandler) buildPatternsInfo(m localization.Messages) string {
	var info string
	for i := 0; i < len(patternEmojis); i++ {
		name, desc := h.getLocalizedPattern(i, m)
		info += fmt.Sprintf("%s *%s*\n_%s_\n\n", patternEmojis[i], name, desc)
	}
	return info
}

func (h *GuidedBreathingHandler) showPatternSelection(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	text := m.GuidedIntro + h.buildPatternsInfo(m)

	buttons := h.buildPatternButtons(m)
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

func (h *GuidedBreathingHandler) buildPatternButtons(m localization.Messages) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	for i := 0; i < len(patternEmojis); i++ {
		name, _ := h.getLocalizedPattern(i, m)
		buttons = append(buttons, []telego.InlineKeyboardButton{
			{Text: fmt.Sprintf("%s %s", patternEmojis[i], name), CallbackData: fmt.Sprintf("gbreath_pattern_%d", i)},
		})
	}
	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: m.Back, CallbackData: "gbreath_cancel"},
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

	if !h.runBreathingCycles(ctx, chatID, userID, messageID, pattern, patternIdx) {
		return
	}

	h.sendCompletion(ctx, chatID, userID, messageID, patternIdx)
}

func (h *GuidedBreathingHandler) runBreathingCycles(
	ctx context.Context, chatID, userID int64, messageID int, pattern techniques.BreathingPattern, patternIdx int,
) bool {
	phases := techniques.GetBreathingPhases(pattern)

	for cycle := 1; cycle <= pattern.Cycles; cycle++ {
		for _, phase := range phases {
			if !h.runPhase(ctx, chatID, userID, messageID, pattern, patternIdx, cycle, phase) {
				return false
			}
		}
	}
	return true
}

func (h *GuidedBreathingHandler) runPhase(
	ctx context.Context, chatID, userID int64, messageID int,
	pattern techniques.BreathingPattern, patternIdx, cycle int, phase techniques.BreathingPhase,
) bool {
	totalSeconds := int(phase.Duration.Seconds())
	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Stop, CallbackData: "gbreath_stop"}},
		},
	}

	patternName, _ := h.getLocalizedPattern(patternIdx, m)
	phaseName := h.getLocalizedPhaseName(phase.Name, m)
	patternEmoji := patternEmojis[patternIdx]

	for elapsed := 0; elapsed <= totalSeconds; elapsed++ {
		if ctx.Err() != nil {
			return false
		}

		state, err := h.sessionStorage.GetState(ctx, userID)
		if err != nil || state != session.StateGuidedBreathingRunning {
			return false
		}

		progress := messages.TimerCountdown(elapsed, totalSeconds)
		text := m.FormatGuidedPhase(patternEmoji, patternName, cycle, pattern.Cycles, phase.Emoji, phaseName, progress)

		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        text,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil && !IsMessageNotModifiedError(err) {
			log.Printf("ERROR: edit guided breathing message: %v", err)
		}

		if elapsed < totalSeconds {
			timer := time.NewTimer(1 * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return false
			case <-timer.C:
			}
		}
	}
	return true
}

func (h *GuidedBreathingHandler) sendCompletion(
	ctx context.Context, chatID, userID int64, messageID int, patternIdx int,
) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateGuidedBreathingRunning {
		return
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	patternName, _ := h.getLocalizedPattern(patternIdx, m)

	text := fmt.Sprintf(m.GuidedCompletion, patternName)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.FeelBetter, CallbackData: "gbreath_complete"},
				{Text: m.Repeat, CallbackData: fmt.Sprintf("gbreath_pattern_%d", patternIdx)},
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

	m := h.localizer.Get(h.getLang(ctx, userID))
	text := m.GuidedStopped + h.buildPatternsInfo(m)

	buttons := h.buildPatternButtons(m)
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

	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.GuidedThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		log.Printf("ERROR: edit guided breathing complete: %v", err)
	}
}

func (h *GuidedBreathingHandler) cancelBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		if !HandleEditError(ctx, h.bot, err, chatID, messageID) {
			log.Printf("ERROR: edit guided breathing cancel: %v", err)
		}
	}
}
