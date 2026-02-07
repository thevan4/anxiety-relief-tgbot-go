// Package handlers implements telegram bot message handlers for anxiety relief techniques.
package handlers

import (
	"context"
	"log"
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

// BreathingHandler handles the basic 4-4-6 breathing exercise.
type BreathingHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.SessionManager
	callbackProcessor *CallbackProcessor
}

// NewBreathingHandler creates a new breathing exercise handler.
func NewBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
	callbackProcessor *CallbackProcessor,
) *BreathingHandler {
	return &BreathingHandler{
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
func (h *BreathingHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// HandleMenuSelect handles selection from main menu.
func (h *BreathingHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showBreathingIntro(h.ctx, info.ChatID, info.UserID)
	return nil
}

// HandleCallback handles breathing exercise callbacks.
func (h *BreathingHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.processCallback(info.ChatID, info.UserID, info.MessageID, info.Data)
	return nil
}

func (h *BreathingHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch data {
	case "breathing_complete":
		h.sessionManager.CancelSession(userID)
		h.completeBreathing(h.ctx, chatID, userID, messageID)
	case "breathing_cancel":
		h.sessionManager.CancelSession(userID)
		h.cancelBreathing(h.ctx, chatID, userID, messageID)
	case "breathing_stop", "breathing_repeat":
		h.sessionManager.CancelSession(userID)
		h.stopBreathing(h.ctx, chatID, userID, messageID)
	case "breathing_pause":
		if err := h.sessionStorage.SetState(h.ctx, userID, session.StateBreathingPaused); err != nil {
			log.Printf("ERROR: set state breathing paused: %v", err)
		}
	case "breathing_resume":
		if err := h.sessionStorage.SetState(h.ctx, userID, session.StateBreathingRunning); err != nil {
			log.Printf("ERROR: set state breathing running: %v", err)
		}
	case "breathing_start":
		// Create new session context (cancels previous if any)
		sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
		if err := h.sessionStorage.SetState(h.ctx, userID, session.StateBreathingRunning); err != nil {
			log.Printf("ERROR: set state breathing running: %v", err)
		}
		// Run in goroutine with session context
		go h.runBreathingCycle(sessionCtx, chatID, userID, messageID)
	}
}

func (h *BreathingHandler) showBreathingIntro(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateBreathingActive); err != nil {
		log.Printf("ERROR: set state breathing active: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "breathing_cancel"},
				{Text: m.Start, CallbackData: "breathing_start"},
			},
		},
	}

	// Recreate message to extend its lifetime
	if _, err := RecreateMenuMessage(ctx, h.bot, h.sessionStorage, chatID, userID, m.BreathingIntro, keyboard); err != nil {
		log.Printf("ERROR: recreate breathing intro: %v", err)
	}
}

const breathingStepsPerCycle = 3

func (h *BreathingHandler) runBreathingCycle(ctx context.Context, chatID, userID int64, messageID int) {
	cycles := techniques.GetBreathingCycles()
	for i, cycle := range cycles {
		if !h.runPhaseWithProgress(ctx, chatID, userID, messageID, i, cycle) {
			return
		}
	}
	h.sendBreathingCompletion(ctx, chatID, userID, messageID)
}

func (h *BreathingHandler) runPhaseWithProgress(
	ctx context.Context, chatID, userID int64, messageID int, stepIndex int, cycle techniques.BreathingCycle,
) bool {
	cycleNum := stepIndex/breathingStepsPerCycle + 1
	totalSeconds := int(cycle.Duration.Seconds())

	m := h.localizer.Get(h.getLang(ctx, userID))

	// Get localized phase name
	phaseName := h.getLocalizedPhaseName(cycle.Phase, m)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Pause, CallbackData: "breathing_pause"}},
		},
	}

	for elapsed := 0; elapsed < totalSeconds; elapsed++ {
		if !h.checkRunningOrPause(ctx, chatID, messageID, userID) {
			return false
		}

		progress := messages.TimerCountdown(elapsed, totalSeconds)
		text := m.FormatBreathingPhase(cycleNum, techniques.BreathCycles, phaseName, cycle.Emoji, progress)

		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        text,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil && !IsMessageNotModifiedError(err) {
			log.Printf("ERROR: edit breathing message: %v", err)
		}

		timer := time.NewTimer(1 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false
		case <-timer.C:
		}
	}
	return true
}

func (h *BreathingHandler) checkRunningOrPause(
	ctx context.Context, chatID int64, messageID int, userID int64,
) bool {
	if ctx.Err() != nil {
		return false
	}
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil {
		return false
	}
	if state == session.StateBreathingPaused {
		h.showPauseScreen(ctx, chatID, messageID, userID)
		return WaitForResume(
			ctx, h.sessionStorage, userID,
			session.StateBreathingRunning, session.StateBreathingPaused,
		)
	}
	return state == session.StateBreathingRunning
}

func (h *BreathingHandler) showPauseScreen(ctx context.Context, chatID int64, messageID int, userID int64) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "breathing_stop"},
				{Text: m.Resume, CallbackData: "breathing_resume"},
			},
		},
	}
	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.PauseText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil && !IsMessageNotModifiedError(err) {
		log.Printf("ERROR: edit breathing pause screen: %v", err)
	}
}

// getLocalizedPhaseName returns localized name for breathing phase.
func (h *BreathingHandler) getLocalizedPhaseName(phaseID string, m localization.Messages) string {
	switch phaseID {
	case "inhale":
		return m.BreathingInhale
	case "hold":
		return m.BreathingHold
	case "exhale":
		return m.BreathingExhale
	default:
		return phaseID
	}
}

func (h *BreathingHandler) sendBreathingCompletion(ctx context.Context, chatID, userID int64, messageID int) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateBreathingRunning {
		return
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Repeat, CallbackData: "breathing_repeat"},
				{Text: m.Done, CallbackData: "breathing_complete"},
			},
		},
	}
	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.BreathingCompletion,
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

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "breathing_cancel"},
				{Text: m.Start, CallbackData: "breathing_start"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.BreathingIntro,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit breathing stop: %v", err)
	}
}

func (h *BreathingHandler) completeBreathing(ctx context.Context, chatID, userID int64, messageID int) {
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
		log.Printf("ERROR: edit breathing complete: %v", err)
	}
}

// getMainMenuInline returns localized main menu keyboard.
func (h *BreathingHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.MenuBreathing, CallbackData: "menu_breathing"}},
			{{Text: m.MenuGrounding, CallbackData: "menu_grounding"}},
			{{Text: m.MenuGuided, CallbackData: "menu_guided"}},
			{{Text: m.MenuPMR, CallbackData: "menu_pmr"}},
			{{Text: m.MenuLang, CallbackData: "menu_lang"}},
		},
	}
}

func (h *BreathingHandler) cancelBreathing(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit breathing cancel: %v", err)
		}
	}
}
