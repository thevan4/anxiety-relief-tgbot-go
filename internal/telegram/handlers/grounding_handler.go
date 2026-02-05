package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
)

// GroundingHandler handles the 5-4-3-2-1 grounding technique.
type GroundingHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.SessionManager
	callbackProcessor *CallbackProcessor
}

// NewGroundingHandler creates a new grounding exercise handler.
func NewGroundingHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
	callbackProcessor *CallbackProcessor,
) *GroundingHandler {
	return &GroundingHandler{
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
func (h *GroundingHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (h *GroundingHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
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
func (h *GroundingHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showGroundingIntro(h.ctx, info.ChatID, info.UserID, info.MessageID)
	return nil
}

// HandleCallback handles grounding exercise callbacks.
func (h *GroundingHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.applyGroundingCallbackAction(h.ctx, info.ChatID, info.UserID, info.MessageID, info.Data)
	return nil
}

func (h *GroundingHandler) applyGroundingCallbackAction(
	ctx context.Context, chatID, userID int64, messageID int, data string,
) {
	switch {
	case strings.HasPrefix(data, "grounding_step_"):
		stepNum := strings.TrimPrefix(data, "grounding_step_")
		if stepNum == "1" || stepNum == "2" || stepNum == "3" || stepNum == "4" || stepNum == "5" {
			h.showGroundingStep(ctx, chatID, userID, messageID, stepNum)
		}
	case data == "grounding_complete":
		h.completeGrounding(ctx, chatID, userID, messageID)
	case data == "grounding_cancel":
		h.cancelGrounding(ctx, chatID, userID, messageID)
	}
}

func (h *GroundingHandler) showGroundingIntro(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateGroundingStep1); err != nil {
		log.Printf("ERROR: set state grounding step1: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "grounding_cancel"},
				{Text: m.Start, CallbackData: "grounding_step_1"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.GroundingIntro,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit grounding intro: %v", err)
	}
}

func getGroundingStepConfig() []struct {
	state    session.State
	nextStep string
} {
	return []struct {
		state    session.State
		nextStep string
	}{
		{session.StateGroundingStep1, "2"},
		{session.StateGroundingStep2, "3"},
		{session.StateGroundingStep3, "4"},
		{session.StateGroundingStep4, "5"},
		{session.StateGroundingStep5, "complete"},
	}
}

// groundingStepEmojis are fixed emojis for each grounding step.
var groundingStepEmojis = []string{"👁️", "🤚", "👂", "👃", "👅"}

func (h *GroundingHandler) showGroundingStep(
	ctx context.Context, chatID, userID int64, messageID int, stepNum string,
) {
	if len(stepNum) != 1 {
		return
	}
	cfg := getGroundingStepConfig()
	idx := int(stepNum[0] - '1')
	if idx < 0 || idx >= len(cfg) {
		return
	}
	stepCfg := cfg[idx]
	if err := h.sessionStorage.SetState(ctx, userID, stepCfg.state); err != nil {
		log.Printf("ERROR: set state grounding: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	nextStep := stepCfg.nextStep

	// Get localized step title and description
	title, desc := h.getLocalizedStep(idx, m)
	emoji := groundingStepEmojis[idx]
	text := fmt.Sprintf("%s *%s*\n\n%s", emoji, title, desc)

	var keyboard *telego.InlineKeyboardMarkup
	if nextStep == "complete" {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: m.Back, CallbackData: "grounding_cancel"},
					{Text: m.Done, CallbackData: "grounding_complete"},
				},
			},
		}
	} else {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: m.Back, CallbackData: "grounding_cancel"},
					{Text: m.Next, CallbackData: "grounding_step_" + nextStep},
				},
			},
		}
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit grounding message: %v", err)
	}
}

// getLocalizedStep returns localized title and description for grounding step.
func (h *GroundingHandler) getLocalizedStep(idx int, m localization.Messages) (string, string) {
	switch idx {
	case 0:
		return m.GroundingStep1Title, m.GroundingStep1Desc
	case 1:
		return m.GroundingStep2Title, m.GroundingStep2Desc
	case 2:
		return m.GroundingStep3Title, m.GroundingStep3Desc
	case 3:
		return m.GroundingStep4Title, m.GroundingStep4Desc
	case 4:
		return m.GroundingStep5Title, m.GroundingStep5Desc
	default:
		return "", ""
	}
}

func (h *GroundingHandler) completeGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.GroundingThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		log.Printf("ERROR: edit grounding complete: %v", err)
	}
}

func (h *GroundingHandler) cancelGrounding(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit grounding cancel: %v", err)
		}
	}
}
