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

func (h *PMRHandler) Handle(ctx *th.Context, update telego.Update) error {
	if update.Message != nil {
		return h.handleMessage(ctx, update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}
	return nil
}

func (h *PMRHandler) handleMessage(_ *th.Context, msg *telego.Message) error {
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
	h.showIntro(h.ctx, msg.Chat.ID, userID)
	return nil
}

func (h *PMRHandler) handleCallback(ctx *th.Context, cb *telego.CallbackQuery) error {
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

func (h *PMRHandler) showIntro(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRActive); err != nil {
		log.Printf("ERROR: set state pmr active: %v", err)
	}

	muscleGroups := techniques.GetMuscleGroups()

	text := fmt.Sprintf(`💪 *Прогрессивная мышечная релаксация*

Техника глубокого расслабления через напряжение и расслабление мышц.

*Как это работает:*
1. Напрягаете группу мышц на 7 секунд
2. Расслабляете на 15 секунд
3. Переходите к следующей группе

*%d групп мышц* — от рук до ног.

Готовы начать?`, len(muscleGroups))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "pmr_start"},
				{Text: "❌ Отмена", CallbackData: "pmr_cancel"},
			},
		},
	}

	removeKeyboard := &telego.ReplyKeyboardRemove{RemoveKeyboard: true}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		text,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard)); err != nil {
		log.Printf("ERROR: send pmr intro: %v", err)
		return
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"_Используйте кнопки выше_",
	).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard)); err != nil {
		log.Printf("ERROR: send pmr use buttons: %v", err)
	}
}

func (h *PMRHandler) startExercise(ctx context.Context, chatID, userID int64, messageID int) {
	h.runMuscleGroup(ctx, chatID, userID, messageID, 0)
}

func (h *PMRHandler) runMuscleGroup(
	ctx context.Context, chatID, userID int64, messageID int, muscleIdx int,
) {
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

	// Tense phase
	tenseText := fmt.Sprintf("🔴 *НАПРЯГИТЕ*\n\n%s\n\n_7 секунд..._", muscle.TenseInstruction)
	if !h.runPhase(ctx, chatID, userID, messageID, muscle, totalGroups, tenseText, muscle.TenseDuration) {
		return
	}

	// Relax phase
	relaxText := fmt.Sprintf("🟢 *РАССЛАБЬТЕ*\n\n%s\n\n_15 секунд..._", muscle.RelaxInstruction)
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

	text := fmt.Sprintf("%s *%s* (%d/%d)\n\n%s",
		muscle.Emoji, muscle.Name, muscle.Number, totalGroups, phaseText)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: "🛑 Стоп", CallbackData: "pmr_stop"}},
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

	text := `✅ *Отлично!*

Вы завершили прогрессивную мышечную релаксацию.

Ваше тело теперь полностью расслаблено. Посидите ещё минуту, наслаждаясь этим состоянием.`

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "✅ Готово", CallbackData: "pmr_complete"},
				{Text: "🔄 Повторить", CallbackData: "pmr_start"},
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
		log.Printf("ERROR: send pmr completion: %v", err)
	}
}

func (h *PMRHandler) stopExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRActive); err != nil {
		log.Printf("ERROR: set state pmr active: %v", err)
	}

	text := `💪 *Прогрессивная мышечная релаксация*

Упражнение остановлено.

Хотите начать сначала?`

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "pmr_start"},
				{Text: "❌ Отмена", CallbackData: "pmr_cancel"},
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
		log.Printf("ERROR: edit pmr stop: %v", err)
	}
}

func (h *PMRHandler) completeExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	text := `✨ *Спасибо за практику!*

Прогрессивная мышечная релаксация снижает мышечное напряжение и уровень стресса.`

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}); err != nil {
		log.Printf("ERROR: edit pmr complete: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send pmr menu: %v", err)
	}
}

func (h *PMRHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      "❌ Упражнение отменено.",
	}); err != nil {
		log.Printf("ERROR: edit pmr cancel: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send pmr menu: %v", err)
	}
}
