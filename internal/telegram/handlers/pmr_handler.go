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
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
)

type PMRHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	localizer      *localization.Localizer
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
	sessionManager *session.SessionManager
}

func NewPMRHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
) *PMRHandler {
	return &PMRHandler{
		ctx:            ctx,
		bot:            bot,
		localizer:      localizer,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
		sessionManager: sessionManager,
	}
}

// getLang returns user's language from session or default.
func (h *PMRHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (h *PMRHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
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
				{Text: m.MenuInfo, CallbackData: "menu_info"},
			},
			{
				{Text: m.MenuLang, CallbackData: "menu_lang"},
			},
		},
	}
}

// HandleMenuSelect handles selection from main menu
func (h *PMRHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
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

	h.showIntro(h.ctx, chatID, userID, messageID)
	return nil
}

// HandleCallback handles PMR exercise callbacks
func (h *PMRHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
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

func (h *PMRHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case data == "pmr_start":
		// Create new session context (cancels previous if any)
		sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
		go h.startExercise(sessionCtx, chatID, userID, messageID)
	case strings.HasPrefix(data, "pmr_muscle_"):
		muscleIdx := strings.TrimPrefix(data, "pmr_muscle_")
		idx, err := strconv.Atoi(muscleIdx)
		if err == nil {
			sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
			go h.runMuscleGroup(sessionCtx, chatID, userID, messageID, idx)
		}
	case data == "pmr_stop":
		h.sessionManager.CancelSession(userID)
		h.stopExercise(h.ctx, chatID, userID, messageID)
	case data == "pmr_complete":
		h.sessionManager.CancelSession(userID)
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case data == "pmr_cancel":
		h.sessionManager.CancelSession(userID)
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *PMRHandler) showIntro(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRActive); err != nil {
		log.Printf("ERROR: set state pmr active: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	muscleGroups := techniques.GetMuscleGroups()
	text := fmt.Sprintf(m.PMRIntro, len(muscleGroups))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "pmr_cancel"},
				{Text: m.Start, CallbackData: "pmr_start"},
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
		log.Printf("ERROR: edit pmr intro: %v", err)
	}
}

func (h *PMRHandler) startExercise(ctx context.Context, chatID, userID int64, messageID int) {
	h.runMuscleGroup(ctx, chatID, userID, messageID, 0)
}

func (h *PMRHandler) runMuscleGroup(ctx context.Context, chatID, userID int64, messageID int, muscleIdx int) {
	muscleGroups := techniques.GetMuscleGroups()
	if muscleIdx < 0 || muscleIdx >= len(muscleGroups) {
		h.sendCompletion(ctx, chatID, userID, messageID)
		return
	}

	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRRunning); err != nil {
		log.Printf("ERROR: set state pmr running: %v", err)
	}

	muscle := muscleGroups[muscleIdx]
	totalGroups := len(muscleGroups)

	// Tense phase with progress
	if !h.runPhaseWithProgress(ctx, chatID, userID, messageID, muscle, totalGroups, true, muscle.TenseDuration) {
		return
	}

	// Relax phase with progress
	if !h.runPhaseWithProgress(ctx, chatID, userID, messageID, muscle, totalGroups, false, muscle.RelaxDuration) {
		return
	}

	h.runMuscleGroup(ctx, chatID, userID, messageID, muscleIdx+1)
}

func (h *PMRHandler) runPhaseWithProgress(
	ctx context.Context, chatID, userID int64, messageID int,
	muscle techniques.MuscleGroup, totalGroups int, isTense bool, duration time.Duration,
) bool {
	totalSeconds := int(duration.Seconds())
	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Stop, CallbackData: "pmr_stop"}},
		},
	}

	for elapsed := 0; elapsed <= totalSeconds; elapsed++ {
		if ctx.Err() != nil {
			return false
		}

		state, err := h.sessionStorage.GetState(ctx, userID)
		if err != nil || state != session.StatePMRRunning {
			return false
		}

		progress := messages.TimerCountdown(elapsed, totalSeconds)
		var phaseText string
		if isTense {
			phaseText = fmt.Sprintf(m.PMRTense, muscle.TenseInstruction) + "\n\n" + progress
		} else {
			phaseText = fmt.Sprintf(m.PMRRelax, muscle.RelaxInstruction) + "\n\n" + progress
		}

		text := fmt.Sprintf("%s *%s* (%d/%d)\n\n%s", muscle.Emoji, muscle.Name, muscle.Number, totalGroups, phaseText)

		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        text,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil && !IsMessageNotModifiedError(err) {
			log.Printf("ERROR: edit pmr phase message: %v", err)
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

func (h *PMRHandler) sendCompletion(ctx context.Context, chatID, userID int64, messageID int) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StatePMRRunning {
		return
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Done, CallbackData: "pmr_complete"},
				{Text: m.Repeat, CallbackData: "pmr_start"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.PMRCompletion,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: send pmr completion: %v", err)
	}
}

func (h *PMRHandler) stopExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRActive); err != nil {
		log.Printf("ERROR: set state pmr active: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "pmr_cancel"},
				{Text: m.Start, CallbackData: "pmr_start"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.PMRStopped,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit pmr stop: %v", err)
	}
}

func (h *PMRHandler) completeExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.PMRThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		log.Printf("ERROR: edit pmr complete: %v", err)
	}
}

func (h *PMRHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit pmr cancel: %v", err)
		}
	}
}
