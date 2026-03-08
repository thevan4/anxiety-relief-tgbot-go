package handlers

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
)

// CallbackProcessor handles callback validation and session management.
// It encapsulates all dependencies needed for proper callback processing.
type CallbackProcessor struct {
	ctx       context.Context
	bot       *telego.Bot
	storage   session.Storage
	localizer *localization.Localizer
}

// NewCallbackProcessor creates a new callback processor with all dependencies.
func NewCallbackProcessor(
	ctx context.Context,
	bot *telego.Bot,
	storage session.Storage,
	localizer *localization.Localizer,
) *CallbackProcessor {
	return &CallbackProcessor{
		ctx:       ctx,
		bot:       bot,
		storage:   storage,
		localizer: localizer,
	}
}

// CallbackInfo contains extracted and validated callback data.
type CallbackInfo struct {
	ChatID    int64
	UserID    int64
	MessageID int
	Data      string
}

// Extract extracts callback info with callback guard.
// Returns nil if callback should not be processed (invalid, too old, stale, or no active menu).
func (cp *CallbackProcessor) Extract(cb telego.CallbackQuery) *CallbackInfo {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		return nil
	}

	info := &CallbackInfo{
		ChatID:    msg.Chat.ID,
		UserID:    cb.From.ID,
		MessageID: msg.MessageID,
		Data:      cb.Data,
	}

	// Callback guard: check if menu exists for this user (before answering so we can show alert once)
	if !cp.hasActiveMenu(info.UserID, info.ChatID, info.MessageID, cb.ID, cb.From.LanguageCode) {
		return nil
	}

	// Answer callback; if too old - stop processing
	if !AnswerCallbackOrDelete(cp.ctx, cp.bot, cb.ID, info.ChatID, info.MessageID) {
		return nil
	}

	return info
}

// hasActiveMenu checks if user has an active menu message (callback guard).
// Returns false and shows alert if no menu or message mismatch.
func (cp *CallbackProcessor) hasActiveMenu(userID, chatID int64, messageID int, callbackID, langCode string) bool {
	storedMsgID, err := cp.storage.GetMenuMessageID(cp.ctx, userID)

	// No menu in Redis — session expired or never existed
	if err != nil || storedMsgID == 0 {
		log.Printf("DEBUG: no menu for user %d, showing session expired", userID)
		cp.showSessionExpired(callbackID, chatID, messageID, userID, langCode)
		return false
	}

	// Message ID mismatch — callback from old/stale message, silently ignore
	if storedMsgID != messageID {
		log.Printf("DEBUG: stale message %d (expected %d) for user %d", messageID, storedMsgID, userID)
		_ = cp.bot.AnswerCallbackQuery(cp.ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: callbackID,
		})
		return false
	}

	return true
}

// showSessionExpired answers the callback and sends a fresh holder message.
func (cp *CallbackProcessor) showSessionExpired(callbackID string, chatID int64, messageID int, userID int64, langCode string) {
	_ = cp.bot.AnswerCallbackQuery(cp.ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
	})

	_ = cp.bot.DeleteMessage(cp.ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	})

	// Don't create new holder if one already exists
	if existingHolderID, err := cp.storage.GetHolderMessageID(cp.ctx, userID); err == nil && existingHolderID != 0 {
		return
	}

	cp.sendHolder(chatID, userID, langCode)
}

// sendHolder sends a new holder (welcome) message with the Start button.
func (cp *CallbackProcessor) sendHolder(chatID, userID int64, langCode string) {
	lang, _ := cp.storage.GetLang(cp.ctx, userID)
	if lang == "" {
		lang = cp.localizer.SupportedLang(langCode)
		_ = cp.storage.SetLang(cp.ctx, userID, lang)
	}
	m := cp.localizer.Get(lang)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Start, CallbackData: "holder_start"}},
		},
	}

	sentMsg, err := cp.bot.SendMessage(cp.ctx, tu.Message(
		tu.ID(chatID),
		m.HolderText,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard))
	if err != nil {
		log.Printf("ERROR: send holder on session expired: %v", err)
		return
	}

	if setErr := cp.storage.SetHolderMessageID(cp.ctx, userID, sentMsg.MessageID); setErr != nil {
		log.Printf("ERROR: save holder message id on session expired: %v", setErr)
	}
}
