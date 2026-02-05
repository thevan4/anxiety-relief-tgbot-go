package handlers

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
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

	// Answer callback; if too old - delete message
	if !AnswerCallbackOrDelete(cp.ctx, cp.bot, cb.ID, info.ChatID, info.MessageID) {
		return nil
	}

	// Callback guard: check if menu exists for this user
	if !cp.hasActiveMenu(info.UserID, info.ChatID, info.MessageID, cb.ID) {
		return nil
	}

	return info
}

// hasActiveMenu checks if user has an active menu message (callback guard).
// Returns false and shows alert if no menu or message mismatch.
func (cp *CallbackProcessor) hasActiveMenu(userID, chatID int64, messageID int, callbackID string) bool {
	storedMsgID, err := cp.storage.GetMenuMessageID(cp.ctx, userID)

	// No menu in Redis — session expired or never existed
	if err != nil || storedMsgID == 0 {
		log.Printf("DEBUG: no menu for user %d, showing session expired", userID)
		cp.showSessionExpired(callbackID, chatID, messageID, userID)
		return false
	}

	// Message ID mismatch — callback from old/stale message
	if storedMsgID != messageID {
		log.Printf("DEBUG: stale message %d (expected %d) for user %d", messageID, storedMsgID, userID)
		cp.showSessionExpired(callbackID, chatID, messageID, userID)
		return false
	}

	return true
}

// showSessionExpired shows alert and deletes stale message.
func (cp *CallbackProcessor) showSessionExpired(callbackID string, chatID int64, messageID int, userID int64) {
	// Get localized message
	lang, _ := cp.storage.GetLang(cp.ctx, userID)
	m := cp.localizer.Get(lang)

	_ = cp.bot.AnswerCallbackQuery(cp.ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
		Text:            m.SessionExpired,
		ShowAlert:       true,
	})
	DeleteMessage(cp.ctx, cp.bot, chatID, messageID)
}
