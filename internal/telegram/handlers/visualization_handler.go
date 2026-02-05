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
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
)

const (
	visualizationStepDuration  = 15 * time.Second
	visualizationIntroDuration = 5 * time.Second
)

// VisualizationHandler handles peaceful scene visualization exercise.
type VisualizationHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.SessionManager
	callbackProcessor *CallbackProcessor
}

// NewVisualizationHandler creates a new visualization exercise handler.
func NewVisualizationHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
	callbackProcessor *CallbackProcessor,
) *VisualizationHandler {
	return &VisualizationHandler{
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
func (h *VisualizationHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (h *VisualizationHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
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
func (h *VisualizationHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showSceneSelectionRecreate(h.ctx, info.ChatID, info.UserID)
	return nil
}

// HandleCallback handles visualization callbacks.
func (h *VisualizationHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.processCallback(info.ChatID, info.UserID, info.MessageID, info.Data)
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

func (h *VisualizationHandler) buildScenesInfo(scenes []localization.VisualizationScene) string {
	var info string
	for _, scene := range scenes {
		info += fmt.Sprintf("%s *%s*\n_%s_\n\n", scene.Emoji, scene.Name, scene.Description)
	}
	return info
}

func (h *VisualizationHandler) showSceneSelectionRecreate(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationSelect); err != nil {
		log.Printf("ERROR: set state visualization select: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	scenes := m.GetVisualizationScenes()
	text := m.VisualizationIntro + h.buildScenesInfo(scenes)

	buttons := h.buildSceneButtons(scenes, m)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	// Recreate message to extend its lifetime
	if _, err := RecreateMenuMessage(ctx, h.bot, h.sessionStorage, chatID, userID, text, keyboard); err != nil {
		log.Printf("ERROR: recreate visualization intro: %v", err)
	}
}

func (h *VisualizationHandler) showSceneSelection(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateVisualizationSelect); err != nil {
		log.Printf("ERROR: set state visualization select: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	scenes := m.GetVisualizationScenes()
	text := m.VisualizationIntro + h.buildScenesInfo(scenes)

	buttons := h.buildSceneButtons(scenes, m)
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
	scenes []localization.VisualizationScene, m localization.Messages,
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
		{Text: m.Back, CallbackData: "visual_cancel"},
	})
	return buttons
}

func (h *VisualizationHandler) startScene(
	ctx context.Context, chatID, userID int64, messageID int, sceneID string,
) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	scene := h.findScene(sceneID, m)
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

func (h *VisualizationHandler) findScene(sceneID string, m localization.Messages) *localization.VisualizationScene {
	scenes := m.GetVisualizationScenes()
	for _, s := range scenes {
		if s.ID == sceneID {
			return &s
		}
	}
	return nil
}

func (h *VisualizationHandler) showSceneIntro(
	ctx context.Context, chatID, userID int64, messageID int, scene *localization.VisualizationScene,
) bool {
	totalSeconds := int(visualizationIntroDuration.Seconds())
	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Stop, CallbackData: "visual_stop"}},
		},
	}

	for elapsed := 0; elapsed <= totalSeconds; elapsed++ {
		if !h.checkVisualizationRunning(ctx, userID) {
			return false
		}

		timer := messages.TimerCountdown(elapsed, totalSeconds)
		introText := m.FormatVisualizationSceneIntro(
			scene.Emoji, scene.Name, scene.Description, scene.Atmosphere, timer,
		)

		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        introText,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil && !IsMessageNotModifiedError(err) {
			log.Printf("ERROR: edit visualization intro: %v", err)
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

func (h *VisualizationHandler) checkVisualizationRunning(ctx context.Context, userID int64) bool {
	if ctx.Err() != nil {
		return false
	}
	state, err := h.sessionStorage.GetState(ctx, userID)
	return err == nil && state == session.StateVisualizationRunning
}

func (h *VisualizationHandler) runSceneSteps(
	ctx context.Context, chatID, userID int64, messageID int, scene *localization.VisualizationScene,
) bool {
	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Stop, CallbackData: "visual_stop"}},
		},
	}

	totalSeconds := int(visualizationStepDuration.Seconds())

	for _, step := range scene.Steps {
		for elapsed := 0; elapsed <= totalSeconds; elapsed++ {
			if !h.checkVisualizationRunning(ctx, userID) {
				return false
			}

			timer := messages.TimerCountdown(elapsed, totalSeconds)
			stepText := m.FormatVisualizationStep(
				scene.Emoji, scene.Name, step.Number, len(scene.Steps),
				step.Instruction, timer,
			)

			if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
				ChatID:      tu.ID(chatID),
				MessageID:   messageID,
				Text:        stepText,
				ParseMode:   "Markdown",
				ReplyMarkup: keyboard,
			}); err != nil && !IsMessageNotModifiedError(err) {
				log.Printf("ERROR: edit visualization step: %v", err)
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

	m := h.localizer.Get(h.getLang(ctx, userID))
	scene := h.findScene(sceneID, m)
	if scene == nil {
		return
	}

	text := m.FormatVisualizationCompletion(scene.Name)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.FeelBetter, CallbackData: "visual_complete"},
				{Text: m.Repeat, CallbackData: fmt.Sprintf("visual_scene_%s", sceneID)},
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

	m := h.localizer.Get(h.getLang(ctx, userID))
	scenes := m.GetVisualizationScenes()
	text := m.VisualizationStopped + h.buildScenesInfo(scenes)

	buttons := h.buildSceneButtons(scenes, m)
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

	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.VisualizationThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		log.Printf("ERROR: edit visualization complete: %v", err)
	}
}

func (h *VisualizationHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit visualization cancel: %v", err)
		}
	}
}
