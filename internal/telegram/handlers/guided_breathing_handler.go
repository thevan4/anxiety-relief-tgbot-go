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

const cbGbreathResume = "gbreath_resume"

// GuidedBreathingHandler handles advanced breathing patterns with different timings.
type GuidedBreathingHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.Manager
	callbackProcessor *CallbackProcessor
}

// NewGuidedBreathingHandler creates a new guided breathing exercise handler.
func NewGuidedBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.Manager,
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
	return GetLang(ctx, h.sessionStorage, userID)
}

// HandleMenuSelect handles selection from main menu.
func (h *GuidedBreathingHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showPatternSelection(h.ctx, info.ChatID, info.UserID)
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
	case strings.HasPrefix(data, "gbreath_pattern_") || strings.HasPrefix(data, "gbreath_start_"):
		h.handlePatternAction(chatID, userID, messageID, data)
	case data == "gbreath_pause" || data == cbGbreathResume:
		h.handlePauseResume(userID, data)
	case data == "gbreath_stop":
		h.sessionManager.CancelSession(userID)
		h.stopBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_complete":
		h.sessionManager.CancelSession(userID)
		h.completeBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_cancel":
		h.sessionManager.CancelSession(userID)
		h.cancelBreathing(h.ctx, chatID, userID, messageID)
	case data == "gbreath_back_select":
		h.showPatternSelectionEdit(h.ctx, chatID, userID, messageID)
	}
}

func (h *GuidedBreathingHandler) handlePatternAction(chatID, userID int64, messageID int, data string) {
	if strings.HasPrefix(data, "gbreath_pattern_") {
		idx, err := strconv.Atoi(strings.TrimPrefix(data, "gbreath_pattern_"))
		if err == nil {
			h.showPatternIntro(h.ctx, chatID, userID, messageID, idx)
		}
		return
	}
	idx, err := strconv.Atoi(strings.TrimPrefix(data, "gbreath_start_"))
	if err == nil {
		sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
		go h.startPattern(sessionCtx, chatID, userID, messageID, idx)
	}
}

func (h *GuidedBreathingHandler) handlePauseResume(userID int64, data string) {
	state := session.StateGuidedBreathingPaused
	if data == cbGbreathResume {
		state = session.StateGuidedBreathingRunning
	}
	if err := h.sessionStorage.SetState(h.ctx, userID, state); err != nil {
		log.Printf("ERROR: set guided breathing state %s: %v", state, err)
	}
}

// patternEmojis are fixed emojis for each breathing pattern.
//
//nolint:gochecknoglobals // Static UI data.
var patternEmojis = []string{"📦", "😴", "⚡", "🚀"}

// getLocalizedPattern returns localized name for pattern by index.
func (h *GuidedBreathingHandler) getLocalizedPattern(idx int, m localization.Messages) string {
	names := []string{
		m.PatternBoxName,
		m.PatternRelaxingName,
		m.PatternEnergizingName,
		m.PatternQuickName,
	}
	if idx < 0 || idx >= len(names) {
		return ""
	}
	return names[idx]
}

// getLocalizedPhaseName returns localized phase name.
func (h *GuidedBreathingHandler) getLocalizedPhaseName(phaseID string, m localization.Messages) string {
	switch phaseID {
	case techniques.PhaseInhale:
		return m.GuidedInhale
	case techniques.PhaseHoldIn:
		return m.GuidedHoldIn
	case techniques.PhaseExhale:
		return m.GuidedExhale
	case techniques.PhaseHoldOut:
		return m.GuidedHoldOut
	default:
		return phaseID
	}
}

// showPatternSelection shows pattern selection screen (from main menu, recreates message).
func (h *GuidedBreathingHandler) showPatternSelection(ctx context.Context, chatID, userID int64) {
	if setErr := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); setErr != nil {
		log.Printf("ERROR: set state guided breathing select: %v", setErr)
	}
	m := h.localizer.Get(h.getLang(ctx, userID))
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: h.buildPatternButtons(m)}
	_, recreateErr := RecreateMenuMessage(
		ctx, h.bot, h.sessionStorage, chatID, userID,
		m.GuidedIntro, keyboard,
	)
	if recreateErr != nil {
		log.Printf("ERROR: recreate guided breathing pattern select: %v", recreateErr)
	}
}

// showPatternSelectionEdit shows pattern selection screen (edit existing message).
func (h *GuidedBreathingHandler) showPatternSelectionEdit(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	buttons := h.buildPatternButtons(m)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.GuidedIntro,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit guided breathing pattern select: %v", err)
	}
}

func (h *GuidedBreathingHandler) buildPatternButtons(m localization.Messages) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	for i := 0; i < len(patternEmojis); i++ {
		name := h.getLocalizedPattern(i, m)
		buttons = append(buttons, []telego.InlineKeyboardButton{
			{Text: fmt.Sprintf("%s %s", patternEmojis[i], name), CallbackData: fmt.Sprintf("gbreath_pattern_%d", i)},
		})
	}
	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: m.Back, CallbackData: "gbreath_cancel"},
	})
	return buttons
}

// getLocalizedPatternIntro returns localized intro text for pattern by index.
func (h *GuidedBreathingHandler) getLocalizedPatternIntro(idx int, m localization.Messages) string {
	intros := []string{
		m.PatternBoxIntro,
		m.PatternRelaxingIntro,
		m.PatternEnergizingIntro,
		m.PatternQuickIntro,
	}
	if idx < 0 || idx >= len(intros) {
		return ""
	}
	return intros[idx]
}

// showPatternIntro shows the detailed intro screen for a specific pattern.
func (h *GuidedBreathingHandler) showPatternIntro(
	ctx context.Context, chatID, userID int64, messageID int, patternIdx int,
) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGuidedBreathingSelect); err != nil {
		log.Printf("ERROR: set state guided breathing select: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	text := h.getLocalizedPatternIntro(patternIdx, m)
	if text == "" {
		return
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "gbreath_back_select"},
				{Text: m.Start, CallbackData: fmt.Sprintf("gbreath_start_%d", patternIdx)},
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
		log.Printf("ERROR: edit guided breathing pattern intro: %v", err)
	}
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
			{{Text: m.Pause, CallbackData: "gbreath_pause"}},
		},
	}

	patternName := h.getLocalizedPattern(patternIdx, m)
	phaseName := h.getLocalizedPhaseName(phase.Phase, m)
	patternEmoji := patternEmojis[patternIdx]

	for elapsed := 0; elapsed < totalSeconds; elapsed++ {
		if !h.checkRunningOrPause(ctx, chatID, messageID, userID) {
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

func (h *GuidedBreathingHandler) checkRunningOrPause(
	ctx context.Context, chatID int64, messageID int, userID int64,
) bool {
	if ctx.Err() != nil {
		return false
	}
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil {
		return false
	}
	if state == session.StateGuidedBreathingPaused {
		h.showPauseScreen(ctx, chatID, messageID, userID)
		return WaitForResume(
			ctx, h.sessionStorage, userID,
			session.StateGuidedBreathingRunning,
			session.StateGuidedBreathingPaused,
		)
	}
	return state == session.StateGuidedBreathingRunning
}

func (h *GuidedBreathingHandler) showPauseScreen(ctx context.Context, chatID int64, messageID int, userID int64) {
	ShowPauseScreen(ctx, h.bot, h.localizer, h.getLang(ctx, userID), chatID, messageID, "gbreath_stop", "gbreath_resume")
}

func (h *GuidedBreathingHandler) sendCompletion(
	ctx context.Context, chatID, userID int64, messageID int, patternIdx int,
) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateGuidedBreathingRunning {
		return
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	patternName := h.getLocalizedPattern(patternIdx, m)

	text := fmt.Sprintf(m.GuidedCompletion, patternName)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Repeat, CallbackData: fmt.Sprintf("gbreath_pattern_%d", patternIdx)},
				{Text: m.Done, CallbackData: "gbreath_complete"},
			},
		},
	}

	if _, editErr := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); editErr != nil {
		log.Printf("ERROR: send guided breathing completion: %v", editErr)
	}
}

func (h *GuidedBreathingHandler) stopBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	h.showPatternSelectionEdit(ctx, chatID, userID, messageID)
}

func (h *GuidedBreathingHandler) completeBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	ClearAndShowMainMenu(
		ctx, h.bot, h.sessionStorage, h.localizer, h.getLang(ctx, userID),
		chatID, userID, messageID, "edit guided breathing complete",
	)
}

func (h *GuidedBreathingHandler) cancelBreathing(ctx context.Context, chatID, userID int64, messageID int) {
	CancelAndShowMainMenu(
		ctx, h.bot, h.sessionStorage, h.localizer, h.getLang(ctx, userID),
		chatID, userID, messageID, "edit guided breathing cancel",
	)
}
