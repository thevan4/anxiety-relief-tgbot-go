package handlers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
)

const pausePollInterval = 500 * time.Millisecond

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

// WaitForResume polls storage until exercise state changes from paused.
// Returns true if resumed (runningState), false otherwise.
func WaitForResume(
	ctx context.Context, storage session.Storage, userID int64,
	runningState, pausedState session.State,
) bool {
	for {
		timer := time.NewTimer(pausePollInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false
		case <-timer.C:
		}
		state, err := storage.GetState(ctx, userID)
		if err != nil {
			return false
		}
		if state == runningState {
			return true
		}
		if state != pausedState {
			return false
		}
	}
}

// RecreateMenuMessage deletes old message and sends a new one.
// Returns new messageID or error. Updates storage with new messageID and menu_created timestamp.
func RecreateMenuMessage(
	ctx context.Context,
	bot *telego.Bot,
	storage session.Storage,
	chatID, userID int64,
	text string,
	keyboard *telego.InlineKeyboardMarkup,
) (int, error) {
	// Get and delete old message
	oldMessageID, _ := storage.GetMenuMessageID(ctx, userID)
	if oldMessageID != 0 {
		_ = bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    telego.ChatID{ID: chatID},
			MessageID: oldMessageID,
		})
	}

	// Send new message
	msg, err := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:      telego.ChatID{ID: chatID},
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return 0, err
	}

	newMessageID := msg.MessageID

	// Update storage with new message ID
	if setErr := storage.SetMenuMessageID(ctx, userID, newMessageID); setErr != nil {
		return 0, setErr
	}

	// Update menu_created for correct cleanup timing
	if setErr := storage.SetMenuCreatedAt(ctx, userID, time.Now()); setErr != nil {
		return 0, setErr
	}

	return newMessageID, nil
}

// GetLang returns user's language from session or default.
func GetLang(ctx context.Context, storage session.Storage, userID int64) string {
	lang, err := storage.GetLang(ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// GetMainMenuInline returns localized main menu keyboard.
func GetMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.MenuBreathing, CallbackData: "menu_breathing"}},
			{{Text: m.MenuGrounding, CallbackData: "menu_grounding"}},
			{{Text: m.MenuGuided, CallbackData: "menu_guided"}},
			{{Text: m.MenuPMR, CallbackData: "menu_pmr"}},
			{{Text: m.MenuLang, CallbackData: "menu_lang"}},
		},
	}
}

// ClearAndShowMainMenu clears user state and edits message to main menu.
func ClearAndShowMainMenu(
	ctx context.Context, bot *telego.Bot, storage session.Storage,
	localizer *localization.Localizer, lang string,
	chatID, userID int64, messageID int, errLog string,
) {
	if clearErr := storage.ClearState(ctx, userID); clearErr != nil {
		log.Printf("ERROR: clear state: %v", clearErr)
	}
	m := localizer.Get(lang)
	if _, editErr := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(m),
	}); editErr != nil {
		log.Printf("ERROR: %s: %v", errLog, editErr)
	}
}

// CancelAndShowMainMenu clears user state and edits message to main menu with error handling.
func CancelAndShowMainMenu(
	ctx context.Context, bot *telego.Bot, storage session.Storage,
	localizer *localization.Localizer, lang string,
	chatID, userID int64, messageID int, errLog string,
) {
	if clearErr := storage.ClearState(ctx, userID); clearErr != nil {
		log.Printf("ERROR: clear state: %v", clearErr)
	}
	m := localizer.Get(lang)
	if _, editErr := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: GetMainMenuInline(m),
	}); editErr != nil {
		if !HandleEditError(ctx, bot, editErr, chatID, messageID) {
			log.Printf("ERROR: %s: %v", errLog, editErr)
		}
	}
}

// ShowPauseScreen shows a pause screen with back and resume buttons.
func ShowPauseScreen(
	ctx context.Context, bot *telego.Bot,
	localizer *localization.Localizer, lang string,
	chatID int64, messageID int,
	stopCB, resumeCB string,
) {
	m := localizer.Get(lang)
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: stopCB},
				{Text: m.Resume, CallbackData: resumeCB},
			},
		},
	}
	if _, editErr := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.PauseText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); editErr != nil && !IsMessageNotModifiedError(editErr) {
		log.Printf("ERROR: edit pause screen: %v", editErr)
	}
}

// ShowIntroRecreate sets state, builds intro with back/start buttons, and recreates the menu message.
func ShowIntroRecreate(
	ctx context.Context, bot *telego.Bot, storage session.Storage,
	localizer *localization.Localizer, lang string,
	chatID, userID int64,
	state session.State, introText string,
	cancelCB, startCB, errLog string,
) {
	if setErr := storage.SetState(ctx, userID, state); setErr != nil {
		log.Printf("ERROR: %s set state: %v", errLog, setErr)
	}
	m := localizer.Get(lang)
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: cancelCB},
				{Text: m.Start, CallbackData: startCB},
			},
		},
	}
	if _, recreateErr := RecreateMenuMessage(ctx, bot, storage, chatID, userID, introText, keyboard); recreateErr != nil {
		log.Printf("ERROR: %s recreate: %v", errLog, recreateErr)
	}
}

// SetStateAndEditIntro sets state and edits message to show intro with back/start buttons.
func SetStateAndEditIntro(
	ctx context.Context, bot *telego.Bot, storage session.Storage,
	chatID, userID int64, messageID int,
	state session.State, introText string,
	backCB, startCB string,
	m localization.Messages, errLog string,
) {
	if setErr := storage.SetState(ctx, userID, state); setErr != nil {
		log.Printf("ERROR: %s set state: %v", errLog, setErr)
	}
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Back, CallbackData: backCB},
				{Text: m.Start, CallbackData: startCB},
			},
		},
	}
	if _, editErr := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        introText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); editErr != nil {
		log.Printf("ERROR: %s edit: %v", errLog, editErr)
	}
}

// SendCompletionScreen checks running state and shows repeat/done buttons.
func SendCompletionScreen(
	ctx context.Context, bot *telego.Bot, storage session.Storage,
	chatID, userID int64, messageID int,
	runningState session.State, completionText string,
	repeatCB, doneCB string,
	m localization.Messages, errLog string,
) {
	state, stateErr := storage.GetState(ctx, userID)
	if stateErr != nil || state != runningState {
		return
	}
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Repeat, CallbackData: repeatCB},
				{Text: m.Done, CallbackData: doneCB},
			},
		},
	}
	if _, editErr := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        completionText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); editErr != nil {
		log.Printf("ERROR: %s: %v", errLog, editErr)
	}
}

// EditCompletionScreen shows repeat/done buttons without state check.
func EditCompletionScreen(
	ctx context.Context, bot *telego.Bot,
	chatID int64, messageID int,
	completionText string,
	repeatCB, doneCB string,
	m localization.Messages, errLog string,
) {
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: m.Repeat, CallbackData: repeatCB},
				{Text: m.Done, CallbackData: doneCB},
			},
		},
	}
	if _, editErr := bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        completionText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); editErr != nil {
		log.Printf("ERROR: %s: %v", errLog, editErr)
	}
}
