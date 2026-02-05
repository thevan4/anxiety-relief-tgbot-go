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

// LangHandler handles language selection.
type LangHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	callbackProcessor *CallbackProcessor
}

// NewLangHandler creates a new LangHandler.
func NewLangHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	callbackProcessor *CallbackProcessor,
) *LangHandler {
	return &LangHandler{
		ctx:               ctx,
		bot:               bot,
		localizer:         localizer,
		statistics:        statistics,
		sessionStorage:    sessionStorage,
		callbackProcessor: callbackProcessor,
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
				{Text: m.MenuVisualization, CallbackData: "menu_visual"},
			},
			{
				{Text: m.MenuLang, CallbackData: "menu_lang"},
			},
		},
	}
}

// HandleMenuSelect handles selection from main menu.
func (h *LangHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showLangSelectionRecreate(h.ctx, info.ChatID, info.UserID)
	return nil
}

// HandleCallback handles language selection callbacks.
func (h *LangHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.processCallback(info.ChatID, info.UserID, info.MessageID, info.Data)
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

func (h *LangHandler) showLangSelectionRecreate(ctx context.Context, chatID, userID int64) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	currentLang := h.getLang(ctx, userID)

	buttons := h.buildLangButtons(currentLang, m)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	// Recreate message to extend its lifetime
	if _, err := RecreateMenuMessage(ctx, h.bot, h.sessionStorage, chatID, userID, m.LangSelectTitle, keyboard); err != nil {
		log.Printf("ERROR: recreate lang selection: %v", err)
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
