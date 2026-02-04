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
)

// LangHandler handles language selection.
type LangHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	localizer      *localization.Localizer
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

// NewLangHandler creates a new LangHandler.
func NewLangHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *LangHandler {
	return &LangHandler{
		ctx:            ctx,
		bot:            bot,
		localizer:      localizer,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

// getLang returns user's language from session or default.
func (h *LangHandler) getLang(ctx context.Context, userID int64) string {
	lang, err := h.sessionStorage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (h *LangHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
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
func (h *LangHandler) HandleMenuSelect(ctx *th.Context, cb telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		return nil
	}
	userID := cb.From.ID
	chatID := msg.Chat.ID
	messageID := msg.MessageID

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

	h.showLangSelection(h.ctx, chatID, userID, messageID)
	return nil
}

// HandleCallback handles language selection callbacks.
func (h *LangHandler) HandleCallback(ctx *th.Context, cb telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		log.Printf("ERROR: callback query message is inaccessible")
		return nil
	}
	chatID := msg.Chat.ID
	userID := cb.From.ID
	messageID := msg.MessageID

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

func (h *LangHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case strings.HasPrefix(data, "lang_set_"):
		lang := strings.TrimPrefix(data, "lang_set_")
		h.setLanguage(h.ctx, chatID, userID, messageID, lang)
	case data == "lang_cancel":
		h.cancelSelection(h.ctx, chatID, userID, messageID)
	}
}

func (h *LangHandler) showLangSelection(ctx context.Context, chatID, userID int64, messageID int) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	currentLang := h.getLang(ctx, userID)

	buttons := h.buildLangButtons(currentLang, m)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.LangSelectTitle,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit lang selection: %v", err)
	}
}

func (h *LangHandler) buildLangButtons(currentLang string, m localization.Messages) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	languages := h.localizer.SupportedLanguages()

	var row []telego.InlineKeyboardButton
	for i, lang := range languages {
		flag := h.localizer.LangFlag(lang)
		name := h.localizer.LangName(lang)
		text := fmt.Sprintf("%s %s", flag, name)
		if lang == currentLang {
			text = fmt.Sprintf("%s %s ✓", flag, name)
		}

		row = append(row, telego.InlineKeyboardButton{
			Text:         text,
			CallbackData: "lang_set_" + lang,
		})

		if len(row) == 2 || i == len(languages)-1 {
			buttons = append(buttons, row)
			row = []telego.InlineKeyboardButton{}
		}
	}

	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: m.Back, CallbackData: "lang_cancel"},
	})

	return buttons
}

func (h *LangHandler) setLanguage(ctx context.Context, chatID, userID int64, messageID int, lang string) {
	normalizedLang := h.localizer.SupportedLang(lang)

	if err := h.sessionStorage.SetLang(ctx, userID, normalizedLang); err != nil {
		log.Printf("ERROR: set language: %v", err)
	}

	m := h.localizer.Get(normalizedLang)

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		log.Printf("ERROR: edit after lang set: %v", err)
	}
}

func (h *LangHandler) cancelSelection(ctx context.Context, chatID, userID int64, messageID int) {
	m := h.localizer.Get(h.getLang(ctx, userID))

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: h.getMainMenuInline(m),
	}); err != nil {
		if !HandleEditError(ctx, h.bot, err, chatID, messageID) {
			log.Printf("ERROR: edit lang cancel: %v", err)
		}
	}
}
