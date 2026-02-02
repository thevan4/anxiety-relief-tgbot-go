package handlers

import (
	"context"
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

type PMRHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewPMRHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *PMRHandler {
	return &PMRHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
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

func (h *PMRHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case data == "pmr_start":
		h.startExercise(h.ctx, chatID, userID, messageID)
	case strings.HasPrefix(data, "pmr_muscle_"):
		muscleIdx := strings.TrimPrefix(data, "pmr_muscle_")
		idx, err := strconv.Atoi(muscleIdx)
		if err == nil {
			h.runMuscleGroup(h.ctx, chatID, userID, messageID, idx)
		}
	case data == "pmr_stop":
		h.stopExercise(h.ctx, chatID, userID, messageID)
	case data == "pmr_complete":
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case data == "pmr_cancel":
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *PMRHandler) showIntro(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRActive); err != nil {
		log.Printf("ERROR: set state pmr active: %v", err)
	}

	muscleGroups := techniques.GetMuscleGroups()
	text := messages.PMRIntro(len(muscleGroups))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.Start, CallbackData: "pmr_start"},
				{Text: messages.Cancel, CallbackData: "pmr_cancel"},
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

	tenseText := messages.PMRTensePhase(muscle.TenseInstruction)
	if !h.runPhase(ctx, chatID, userID, messageID, muscle, totalGroups, tenseText, muscle.TenseDuration) {
		return
	}

	relaxText := messages.PMRRelaxPhase(muscle.RelaxInstruction)
	if !h.runPhase(ctx, chatID, userID, messageID, muscle, totalGroups, relaxText, muscle.RelaxDuration) {
		return
	}

	h.runMuscleGroup(ctx, chatID, userID, messageID, muscleIdx+1)
}

func (h *PMRHandler) runPhase(
	ctx context.Context, chatID, userID int64, messageID int,
	muscle techniques.MuscleGroup, totalGroups int, phaseText string, duration time.Duration,
) bool {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StatePMRRunning {
		return false
	}

	text := messages.PMRMuscleStep(muscle.Emoji, muscle.Name, muscle.Number, totalGroups, phaseText)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: messages.Stop, CallbackData: "pmr_stop"}},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit pmr phase message: %v", err)
	}

	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		timer.Stop()
		return false
	case <-timer.C:
		return true
	}
}

func (h *PMRHandler) sendCompletion(ctx context.Context, chatID, userID int64, messageID int) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StatePMRRunning {
		return
	}

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.Done, CallbackData: "pmr_complete"},
				{Text: messages.Repeat, CallbackData: "pmr_start"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.PMRCompletion,
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

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.Start, CallbackData: "pmr_start"},
				{Text: messages.Cancel, CallbackData: "pmr_cancel"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.PMRStopped,
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

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.PMRThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit pmr complete: %v", err)
	}
}

func (h *PMRHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
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
		log.Printf("ERROR: edit pmr cancel: %v", err)
	}
}
