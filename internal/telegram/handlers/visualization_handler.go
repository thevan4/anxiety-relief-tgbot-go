package handlers

import (
	"context"
	"fmt"
	"log"
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

const (
	visualizationStepDuration  = 15 * time.Second
	visualizationIntroDuration = 5 * time.Second
)

type VisualizationHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewVisualizationHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *VisualizationHandler {
	return &VisualizationHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

func (h *VisualizationHandler) Handle(ctx *th.Context, update telego.Update) error {
	if update.Message != nil {
		return h.handleMessage(ctx, update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}
	return nil
}

func (h *VisualizationHandler) handleMessage(_ *th.Context, msg *telego.Message) error {
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
	h.showSceneSelection(h.ctx, msg.Chat.ID, userID)
	return nil
}

func (h *VisualizationHandler) handleCallback(ctx *th.Context, cb *telego.CallbackQuery) error {
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

func (h *VisualizationHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case strings.HasPrefix(data, "visual_scene_"):
		sceneID := strings.TrimPrefix(data, "visual_scene_")
		h.startScene(h.ctx, chatID, userID, messageID, sceneID)
	case data == "visual_stop":
		h.stopExercise(h.ctx, chatID, userID, messageID)
	case data == "visual_complete":
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case data == "visual_cancel":
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *VisualizationHandler) showSceneSelection(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationSelect); err != nil {
		log.Printf("ERROR: set state visualization select: %v", err)
	}

	scenes := techniques.GetVisualizationScenes()

	text := "🌅 *Мирная визуализация*\n\nТехника расслабления через воображение спокойных мест.\n\n"
	text += "*Выберите сцену:*\n\n"
	for _, scene := range scenes {
		text += fmt.Sprintf("%s *%s*\n_%s_\n\n", scene.Emoji, scene.Name, scene.Description)
	}

	buttons := h.buildSceneButtons(scenes)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	removeKeyboard := &telego.ReplyKeyboardRemove{RemoveKeyboard: true}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		text,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard)); err != nil {
		log.Printf("ERROR: send visualization intro: %v", err)
		return
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"_Используйте кнопки выше_",
	).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard)); err != nil {
		log.Printf("ERROR: send visualization use buttons: %v", err)
	}
}

func (h *VisualizationHandler) buildSceneButtons(
	scenes []techniques.VisualizationScene,
) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	for _, scene := range scenes {
		buttons = append(buttons, []telego.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s %s", scene.Emoji, scene.Name),
				CallbackData: fmt.Sprintf("visual_scene_%s", scene.ID),
			},
		})
	}
	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: "❌ Отмена", CallbackData: "visual_cancel"},
	})
	return buttons
}

func (h *VisualizationHandler) startScene(
	ctx context.Context, chatID, userID int64, messageID int, sceneID string,
) {
	scene := h.findScene(sceneID)
	if scene == nil {
		return
	}

	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationRunning); err != nil {
		log.Printf("ERROR: set state visualization running: %v", err)
	}

	if !h.showSceneIntro(ctx, chatID, messageID, scene) {
		return
	}

	if !h.runSceneSteps(ctx, chatID, userID, messageID, scene) {
		return
	}

	h.sendCompletion(ctx, chatID, userID, messageID, sceneID)
}

func (h *VisualizationHandler) findScene(sceneID string) *techniques.VisualizationScene {
	scenes := techniques.GetVisualizationScenes()
	for _, s := range scenes {
		if s.ID == sceneID {
			return &s
		}
	}
	return nil
}

func (h *VisualizationHandler) showSceneIntro(
	ctx context.Context, chatID int64, messageID int, scene *techniques.VisualizationScene,
) bool {
	introText := fmt.Sprintf(`%s *%s*

_%s_

🌿 *Атмосфера:* %s

Закройте глаза и погрузитесь в эту сцену...`, scene.Emoji, scene.Name, scene.Description, scene.Atmosphere)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: "🛑 Стоп", CallbackData: "visual_stop"}},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        introText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit visualization intro: %v", err)
	}

	timer := time.NewTimer(visualizationIntroDuration)
	select {
	case <-ctx.Done():
		timer.Stop()
		return false
	case <-timer.C:
		return true
	}
}

func (h *VisualizationHandler) runSceneSteps(
	ctx context.Context, chatID, userID int64, messageID int, scene *techniques.VisualizationScene,
) bool {
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: "🛑 Стоп", CallbackData: "visual_stop"}},
		},
	}

	for _, step := range scene.Steps {
		state, err := h.sessionStorage.GetState(ctx, userID)
		if err != nil || state != session.StateVisualizationRunning {
			return false
		}

		stepText := fmt.Sprintf(`%s *%s*

*Шаг %d из %d*

%s`, scene.Emoji, scene.Name, step.Number, len(scene.Steps), step.Instruction)

		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        stepText,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil {
			log.Printf("ERROR: edit visualization step: %v", err)
		}

		timer := time.NewTimer(visualizationStepDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false
		case <-timer.C:
		}
	}
	return true
}

func (h *VisualizationHandler) sendCompletion(
	ctx context.Context, chatID, userID int64, messageID int, sceneID string,
) {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateVisualizationRunning {
		return
	}

	scene := h.findScene(sceneID)
	if scene == nil {
		return
	}

	text := fmt.Sprintf(`✅ *Отлично!*

Вы завершили визуализацию "%s".

Медленно возвращайтесь в реальность. Пошевелите пальцами, глубоко вдохните и откройте глаза.

Как вы себя чувствуете?`, scene.Name)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "✅ Лучше", CallbackData: "visual_complete"},
				{Text: "🔄 Повторить", CallbackData: fmt.Sprintf("visual_scene_%s", sceneID)},
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
		log.Printf("ERROR: send visualization completion: %v", err)
	}
}

func (h *VisualizationHandler) stopExercise(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationSelect); err != nil {
		log.Printf("ERROR: set state visualization select: %v", err)
	}

	scenes := techniques.GetVisualizationScenes()

	text := "🌅 *Мирная визуализация*\n\nУпражнение остановлено. Выберите другую сцену:\n\n"
	for _, scene := range scenes {
		text += fmt.Sprintf("%s *%s*\n", scene.Emoji, scene.Name)
	}

	buttons := h.buildSceneButtons(scenes)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit visualization stop: %v", err)
	}
}

func (h *VisualizationHandler) completeExercise(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	text := `✨ *Спасибо за практику!*

Визуализация — мощная техника для снижения стресса и тревожности. Регулярная практика усиливает эффект.`

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}); err != nil {
		log.Printf("ERROR: edit visualization complete: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send visualization menu: %v", err)
	}
}

func (h *VisualizationHandler) cancelExercise(
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
		log.Printf("ERROR: edit visualization cancel: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send visualization menu: %v", err)
	}
}
