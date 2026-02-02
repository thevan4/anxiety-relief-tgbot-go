package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
)

type ThoughtLabelingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
	sessionManager *session.SessionManager
}

func NewThoughtLabelingHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.SessionManager,
) *ThoughtLabelingHandler {
	return &ThoughtLabelingHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
		sessionManager: sessionManager,
	}
}

// HandleMenuSelect handles selection from main menu
func (h *ThoughtLabelingHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
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

// HandleCallback handles thought labeling callbacks
func (h *ThoughtLabelingHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
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

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.Start, CallbackData: "thought_start"},
				{Text: messages.Cancel, CallbackData: "thought_cancel"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.ThoughtLabelingIntro,
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

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: messages.Cancel, CallbackData: "thought_cancel"}},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.ThoughtLabelingPrompt,
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

// ProcessThoughtInput processes the user's thought text input
func (h *ThoughtLabelingHandler) ProcessThoughtInput(chatID, userID int64, messageID int, thought string) {
	h.showCategories(h.ctx, chatID, userID, messageID, thought)
}

func (h *ThoughtLabelingHandler) showCategories(ctx context.Context, chatID, userID int64, messageID int, thought string) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingCategory); err != nil {
		log.Printf("ERROR: set state thought labeling category: %v", err)
	}

	categories := techniques.GetThoughtCategories()
	text := messages.ThoughtLabelingCategories(thought)

	buttons := h.buildCategoryButtons(categories)
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
	categories []techniques.ThoughtCategory,
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
		{Text: messages.Cancel, CallbackData: "thought_cancel"},
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

	text := messages.ThoughtLabelingCategoryInfo(category.Emoji, category.Name, category.Description, category.Example)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: messages.AnotherThought, CallbackData: "thought_another"},
				{Text: messages.Done, CallbackData: "thought_complete"},
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

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        messages.ThoughtLabelingThanks,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(),
	}); err != nil {
		log.Printf("ERROR: edit thought complete: %v", err)
	}
}

func (h *ThoughtLabelingHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
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
			log.Printf("ERROR: edit thought cancel: %v", err)
		}
	}
}
