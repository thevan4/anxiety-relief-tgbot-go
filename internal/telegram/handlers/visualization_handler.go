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
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
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
	sessionManager *session.SessionManager
}

func NewVisualizationHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
) *VisualizationHandler {
	return &VisualizationHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
		sessionManager: sessionManager,
	}
}

// HandleMenuSelect handles selection from main menu
func (h *VisualizationHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
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

	h.showSceneSelection(h.ctx, chatID, userID, messageID)
	return nil
}

// HandleCallback handles visualization callbacks
func (h *VisualizationHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
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

func (h *VisualizationHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case strings.HasPrefix(data, "visual_scene_"):
		sceneID := strings.TrimPrefix(data, "visual_scene_")
		// Create new session context (cancels previous if any)
		sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
		go h.startScene(sessionCtx, chatID, userID, messageID, sceneID)
	case data == "visual_stop":
		h.sessionManager.CancelSession(userID)
		h.stopExercise(h.ctx, chatID, userID, messageID)
	case data == "visual_complete":
		h.sessionManager.CancelSession(userID)
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case data == "visual_cancel":
		h.sessionManager.CancelSession(userID)
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *VisualizationHandler) buildScenesInfo(scenes []techniques.VisualizationScene) string {
	var info string
	for _, scene := range scenes {
		info += fmt.Sprintf("%s *%s*\n_%s_\n\n", scene.Emoji, scene.Name, scene.Description)
	}
	return info
}

func (h *VisualizationHandler) showSceneSelection(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationSelect); err != nil {
		log.Printf("ERROR: set state visualization select: %v", err)
	}

	scenes := techniques.GetVisualizationScenes()
	text := messages.VisualizationIntro(h.buildScenesInfo(scenes))

	buttons := h.buildSceneButtons(scenes)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit visualization intro: %v", err)
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
		{Text: messages.Cancel, CallbackData: "visual_cancel"},
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

	if !h.showSceneIntro(ctx, chatID, userID, messageID, scene) {
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
	ctx context.Context, chatID, userID int64, messageID int, scene *techniques.VisualizationScene,
) bool {
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil || state != session.StateVisualizationRunning {
		return false
	}

	introText := messages.VisualizationSceneIntro(scene.Emoji, scene.Name, scene.Description, scene.Atmosphere)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: messages.Stop, CallbackData: "visual_stop"}},
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
			{{Text: messages.Stop, CallbackData: "visual_stop"}},
		},
	}

	for _, step := range scene.Steps {
		state, err := h.sessionStorage.GetState(ctx, userID)
		if err != nil || state != session.StateVisualizationRunning {
			return false
		}

		stepText := messages.VisualizationStep(scene.Emoji, scene.Name, step.Number, len(scene.Steps), step.Instruction)

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

	text := messages.VisualizationCompletion(scene.Name)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.FeelBetter, CallbackData: "visual_complete"},
				{Text: messages.Repeat, CallbackData: fmt.Sprintf("visual_scene_%s", sceneID)},
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

func (h *VisualizationHandler) stopExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationSelect); err != nil {
		log.Printf("ERROR: set state visualization select: %v", err)
	}

	scenes := techniques.GetVisualizationScenes()
	text := messages.VisualizationStoppedIntro(h.buildScenesInfo(scenes))

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

func (h *VisualizationHandler) completeExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.VisualizationThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit visualization complete: %v", err)
	}
}

func (h *VisualizationHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit visualization cancel: %v", err)
		}
	}
}
