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

const cbGroundingComplete = "grounding_complete"

// GroundingHandler handles the 5-4-3-2-1 grounding technique.
type GroundingHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.Manager
	callbackProcessor *CallbackProcessor
}

// NewGroundingHandler creates a new grounding exercise handler.
func NewGroundingHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.Manager,
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
	return GetLang(ctx, h.sessionStorage, userID)
}

// HandleMenuSelect handles selection from main menu.
func (h *GroundingHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showGroundingIntro(h.ctx, info.ChatID, info.UserID)
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
		h.showGroundingStep(ctx, chatID, userID, messageID, stepNum)
	case data == "grounding_back_intro":
		h.showGroundingIntroEdit(ctx, chatID, userID, messageID)
	case data == cbGroundingComplete:
		h.showGroundingCompletion(ctx, chatID, userID, messageID)
	case data == "grounding_repeat":
		h.showGroundingIntroEdit(ctx, chatID, userID, messageID)
	case data == "grounding_done":
		h.completeGrounding(ctx, chatID, userID, messageID)
	case data == "grounding_cancel":
		h.cancelGrounding(ctx, chatID, userID, messageID)
	}
}

func (h *GroundingHandler) showGroundingIntro(ctx context.Context, chatID, userID int64) {
	lang := h.getLang(ctx, userID)
	m := h.localizer.Get(lang)
	ShowIntroRecreate(
		ctx, h.bot, h.sessionStorage, h.localizer, lang,
		chatID, userID, session.StateGroundingStep1, m.GroundingIntro,
		"grounding_cancel", "grounding_step_1", "grounding intro",
	)
}

func (h *GroundingHandler) showGroundingIntroEdit(ctx context.Context, chatID, userID int64, messageID int) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	SetStateAndEditIntro(
		ctx, h.bot, h.sessionStorage, chatID, userID, messageID,
		session.StateGroundingStep1, m.GroundingIntro,
		"grounding_cancel", "grounding_step_1", m, "grounding intro edit",
	)
}

func getGroundingStepConfig() []struct {
	state    session.State
	nextStep string
	backCB   string
} {
	return []struct {
		state    session.State
		nextStep string
		backCB   string
	}{
		{session.StateGroundingStep1, "2", "grounding_back_intro"},
		{session.StateGroundingStep2, "3", "grounding_step_1"},
		{session.StateGroundingStep3, "4", "grounding_step_2"},
		{session.StateGroundingStep4, "5", "grounding_step_3"},
		{session.StateGroundingStep5, "complete", "grounding_step_4"},
	}
}

// groundingStepEmojis are fixed emojis for each grounding step.
//
//nolint:gochecknoglobals // Static UI data.
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
	backCB := stepCfg.backCB

	// Get localized step title and description
	title, desc := h.getLocalizedStep(idx, m)
	emoji := groundingStepEmojis[idx]
	text := fmt.Sprintf("%s *%s*\n\n%s", emoji, title, desc)

	var keyboard *telego.InlineKeyboardMarkup
	var nextCB string
	if nextStep == "complete" {
		nextCB = cbGroundingComplete
	} else {
		nextCB = "grounding_step_" + nextStep
	}
	keyboard = &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: backCB},
				{Text: m.Next, CallbackData: nextCB},
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
		log.Printf("ERROR: edit grounding message: %v", err)
	}
}

// getLocalizedStep returns localized title and description for grounding step.
func (h *GroundingHandler) getLocalizedStep(idx int, m localization.Messages) (string, string) {
	steps := []struct {
		title string
		desc  string
	}{
		{m.GroundingStep1Title, m.GroundingStep1Desc},
		{m.GroundingStep2Title, m.GroundingStep2Desc},
		{m.GroundingStep3Title, m.GroundingStep3Desc},
		{m.GroundingStep4Title, m.GroundingStep4Desc},
		{m.GroundingStep5Title, m.GroundingStep5Desc},
	}
	if idx < 0 || idx >= len(steps) {
		return "", ""
	}
	return steps[idx].title, steps[idx].desc
}

func (h *GroundingHandler) showGroundingCompletion(ctx context.Context, chatID, userID int64, messageID int) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	EditCompletionScreen(
		ctx, h.bot, chatID, messageID, m.GroundingCompletion,
		"grounding_repeat", "grounding_done", m, "edit grounding completion",
	)
}

func (h *GroundingHandler) completeGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	ClearAndShowMainMenu(
		ctx, h.bot, h.sessionStorage, h.localizer, h.getLang(ctx, userID),
		chatID, userID, messageID, "edit grounding complete",
	)
}

func (h *GroundingHandler) cancelGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	CancelAndShowMainMenu(
		ctx, h.bot, h.sessionStorage, h.localizer, h.getLang(ctx, userID),
		chatID, userID, messageID, "edit grounding cancel",
	)
}
