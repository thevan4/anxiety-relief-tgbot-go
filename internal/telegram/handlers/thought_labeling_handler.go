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
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
)

// ThoughtLabelingHandler handles thought labeling and categorization exercise.
type ThoughtLabelingHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	rateLimiter       rate_limiter.Limiter
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.SessionManager
	callbackProcessor *CallbackProcessor
}

// NewThoughtLabelingHandler creates a new thought labeling exercise handler.
func NewThoughtLabelingHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
	callbackProcessor *CallbackProcessor,
) *ThoughtLabelingHandler {
	return &ThoughtLabelingHandler{
		ctx:               ctx,
		bot:               bot,
		localizer:         localizer,
		rateLimiter:       rateLimiter,
		statistics:        statistics,
		sessionStorage:    sessionStorage,
		sessionManager:    sessionManager,
		callbackProcessor: callbackProcessor,
	}
}

// getLang returns user's language from session or default.
func (h *ThoughtLabelingHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (h *ThoughtLabelingHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
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

// HandleMenuSelect handles selection from main menu.
func (h *ThoughtLabelingHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.rateLimiter.WaitAndGo(h.ctx, info.UserID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showIntro(h.ctx, info.ChatID, info.UserID, info.MessageID)
	return nil
}

// HandleCallback handles thought labeling callbacks.
func (h *ThoughtLabelingHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.rateLimiter.WaitAndGo(h.ctx, info.UserID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}

	h.processCallback(info.ChatID, info.UserID, info.MessageID, info.Data)
	return nil
}

func (h *ThoughtLabelingHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case data == "thought_start":
		h.promptForThought(h.ctx, chatID, userID, messageID)
	case strings.HasPrefix(data, "thought_cat_"):
		categoryID := strings.TrimPrefix(data, "thought_cat_")
		h.showCategoryInfo(h.ctx, chatID, userID, messageID, categoryID)
	case data == "thought_another":
		h.promptForThought(h.ctx, chatID, userID, messageID)
	case data == "thought_complete":
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case data == "thought_cancel":
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *ThoughtLabelingHandler) showIntro(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingActive); err != nil {
		log.Printf("ERROR: set state thought labeling active: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: "thought_cancel"},
				{Text: m.Start, CallbackData: "thought_start"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.ThoughtIntro,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit thought labeling intro: %v", err)
	}

	// Save message ID for later editing when user sends thought
	if err := h.sessionStorage.SetMessageID(ctx, userID, messageID); err != nil {
		log.Printf("ERROR: save message id: %v", err)
	}
}

func (h *ThoughtLabelingHandler) promptForThought(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingInput); err != nil {
		log.Printf("ERROR: set state thought labeling input: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Back, CallbackData: "thought_cancel"}},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.ThoughtPrompt,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit thought prompt: %v", err)
	}

	// Save message ID for later editing when user sends thought
	if err := h.sessionStorage.SetMessageID(ctx, userID, messageID); err != nil {
		log.Printf("ERROR: save message id: %v", err)
	}
}

// ProcessThoughtInput processes the user's thought text input.
func (h *ThoughtLabelingHandler) ProcessThoughtInput(chatID, userID int64, messageID int, thought string) {
	h.showCategories(h.ctx, chatID, userID, messageID, thought)
}

func (h *ThoughtLabelingHandler) showCategories(
	ctx context.Context, chatID, userID int64, messageID int, thought string,
) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingCategory); err != nil {
		log.Printf("ERROR: set state thought labeling category: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	categories := techniques.GetThoughtCategories()
	text := fmt.Sprintf(m.ThoughtCategories, thought)

	buttons := h.buildCategoryButtons(categories, m)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit thought categories: %v", err)
	}
}

func (h *ThoughtLabelingHandler) buildCategoryButtons(
	categories []techniques.ThoughtCategory, m localization.Messages,
) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	row := []telego.InlineKeyboardButton{}
	for i, cat := range categories {
		row = append(row, telego.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s %s", cat.Emoji, cat.Name),
			CallbackData: fmt.Sprintf("thought_cat_%s", cat.ID),
		})
		if len(row) == 2 || i == len(categories)-1 {
			buttons = append(buttons, row)
			row = []telego.InlineKeyboardButton{}
		}
	}
	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: m.Back, CallbackData: "thought_cancel"},
	})
	return buttons
}

func (h *ThoughtLabelingHandler) showCategoryInfo(
	ctx context.Context, chatID, userID int64, messageID int, categoryID string,
) {
	categories := techniques.GetThoughtCategories()
	var category *techniques.ThoughtCategory
	for _, cat := range categories {
		if cat.ID == categoryID {
			category = &cat
			break
		}
	}
	if category == nil {
		return
	}

	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingActive); err != nil {
		log.Printf("ERROR: set state thought labeling active: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))
	text := fmt.Sprintf(m.ThoughtResult, category.Emoji, category.Name, category.Description, category.Example) +
		"\n\n" + m.ThoughtCompletion

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Repeat, CallbackData: "thought_another"},
				{Text: m.Done, CallbackData: "thought_complete"},
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
		log.Printf("ERROR: edit thought category info: %v", err)
	}
}

func (h *ThoughtLabelingHandler) completeExercise(ctx context.Context, chatID, userID int64, messageID int) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.ThoughtThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		log.Printf("ERROR: edit thought complete: %v", err)
	}
}

func (h *ThoughtLabelingHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit thought cancel: %v", err)
		}
	}
}
