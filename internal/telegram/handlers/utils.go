package handlers

import (
	"context"
	"strings"

	"github.com/mymmrac/telego"
)

// AnswerCallbackOrDelete answers callback query. If "too old" error occurs, deletes the message.
// Returns true if callback was answered successfully, false if message was deleted (stop processing).
func AnswerCallbackOrDelete(ctx context.Context, bot *telego.Bot, callbackID string, chatID int64, messageID int) bool {
	err := bot.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
	})
	if err == nil {
		return true
	}
	if IsCallbackTooOldError(err) {
		DeleteMessage(ctx, bot, chatID, messageID)
		return false
	}
	return true
}

// IsCallbackTooOldError checks if the error is "query is too old" error.
func IsCallbackTooOldError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "query is too old") ||
		strings.Contains(errStr, "query ID is invalid")
}

// IsMessageNotModifiedError checks if the error is "message is not modified" error.
func IsMessageNotModifiedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "message is not modified")
}

// HandleEditError checks edit error and deletes message if it's "not modified" error.
// Returns true if error was handled (message deleted), false otherwise.
func HandleEditError(ctx context.Context, bot *telego.Bot, err error, chatID int64, messageID int) bool {
	if err == nil {
		return false
	}
	if IsMessageNotModifiedError(err) {
		DeleteMessage(ctx, bot, chatID, messageID)
		return true
	}
	return false
}

// DeleteMessage deletes a message from chat.
func DeleteMessage(ctx context.Context, bot *telego.Bot, chatID int64, messageID int) {
	_ = bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    telego.ChatID{ID: chatID},
		MessageID: messageID,
	})
}
