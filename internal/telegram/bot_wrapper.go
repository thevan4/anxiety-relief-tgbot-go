// Package telegram implements the Telegram bot handlers and routing.
package telegram

import (
	"context"
	"log"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/handlers"
)

// BotHandler coordinates all bot message handling and routing.
type BotHandler struct {
	ctx               context.Context
	cancelFunc        context.CancelFunc
	bot               *telego.Bot
	handler           *th.BotHandler
	localizer         *localization.Localizer
	rateLimiter       rate_limiter.Limiter
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.SessionManager
	callbackProcessor *handlers.CallbackProcessor
}

// MustNewBotHandler creates a new bot handler or panics on error.
func MustNewBotHandler(
	ctx context.Context,
	token string,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	options ...telego.BotOption,
) *BotHandler {
	bot, err := telego.NewBot(token, options...)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	pollingCtx, cancel := context.WithCancel(ctx)

	updates, err := bot.UpdatesViaLongPolling(pollingCtx, nil)
	if err != nil {
		cancel()
		log.Fatalf("failed to start long polling: %v", err)
	}

	botHandler, err := th.NewBotHandler(bot, updates)
	if err != nil {
		cancel()
		log.Fatalf("failed to create bot handler: %v", err)
	}

	localizer := localization.NewLocalizer()

	bh := &BotHandler{
		ctx:            ctx,
		cancelFunc:     cancel,
		bot:            bot,
		handler:        botHandler,
		localizer:      localizer,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
		sessionManager: session.NewSessionManager(),
	}

	bh.callbackProcessor = handlers.NewCallbackProcessor(
		ctx, bot, sessionStorage, localizer,
	)

	bh.registerHandlers()

	return bh
}

func (bh *BotHandler) registerHandlers() {
	bh.registerStartHandler()
	bh.registerMenuCallbackHandler()
	bh.registerTechniqueHandlers()
	bh.registerCatchAllHandler() // Must be last!
}

// getLang returns user's language from session or default.
func (bh *BotHandler) getLang(userID int64) string {
	lang, err := bh.sessionStorage.GetLang(bh.ctx, userID)
	if err != nil || lang == "" {
		return localization.DefaultLang
	}
	return lang
}

// getMainMenuInline returns localized main menu keyboard.
func (bh *BotHandler) getMainMenuInline(m localization.Messages) *telego.InlineKeyboardMarkup {
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

// deleteOldMessagesResult contains results of deleting old messages.
type deleteOldMessagesResult struct {
	holderExisted bool
	holderDeleted bool
	menuExisted   bool
	menuDeleted   bool
}

// tryDeleteOldMessages attempts to delete old holder and menu messages.
// Returns struct with deletion results for different handling.
func (bh *BotHandler) tryDeleteOldMessages(chatID, userID int64) deleteOldMessagesResult {
	result := deleteOldMessagesResult{}

	if oldHolderID, err := bh.sessionStorage.GetHolderMessageID(bh.ctx, userID); err == nil && oldHolderID != 0 {
		result.holderExisted = true
		if err := bh.bot.DeleteMessage(bh.ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(chatID),
			MessageID: oldHolderID,
		}); err != nil {
			log.Printf("DEBUG: could not delete old holder %d: %v", oldHolderID, err)
		} else {
			result.holderDeleted = true
		}
	}

	if oldMenuID, err := bh.sessionStorage.GetMenuMessageID(bh.ctx, userID); err == nil && oldMenuID != 0 {
		result.menuExisted = true
		if err := bh.bot.DeleteMessage(bh.ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(chatID),
			MessageID: oldMenuID,
		}); err != nil {
			log.Printf("DEBUG: could not delete old menu %d: %v", oldMenuID, err)
		} else {
			result.menuDeleted = true
		}
	}

	return result
}

func (bh *BotHandler) registerStartHandler() {
	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if message.From == nil {
			return nil
		}
		userID := message.From.ID
		chatID := message.Chat.ID

		bh.rateLimiter.WaitAndGo(bh.ctx, userID)
		if bh.ctx.Err() != nil {
			return bh.ctx.Err()
		}

		// 1. Cancel active exercise
		bh.sessionManager.CancelSession(userID)

		// 2-3. Try to delete old holder and menu
		delResult := bh.tryDeleteOldMessages(chatID, userID)

		// 4. Clear session data
		if err := bh.sessionStorage.ClearSession(bh.ctx, userID); err != nil {
			log.Printf("ERROR: clear session on start: %v", err)
		}

		// 5. Clear cleanup queue and retry
		_ = bh.sessionStorage.RemoveFromCleanupQueue(bh.ctx, userID)
		_ = bh.sessionStorage.ClearCleanupRetry(bh.ctx, userID)

		bh.statistics.IncreaseRequestsStatisticForUser(
			userID, message.From.Username, message.From.IsPremium, message.From.IsBot,
		)

		// Decision table:
		// holderDeleted=true  → create new holder
		// holderDeleted=false → don't create (old holder still works)
		// Note: "deleted=true" means message was deleted OR didn't exist in Redis

		if delResult.holderExisted && !delResult.holderDeleted {
			// Holder existed but couldn't be deleted (>48h) → don't create new
			// User still sees old holder with "Start" button
			log.Printf("DEBUG: old holder not deleted, skipping new holder for user %d", userID)
			return nil
		}

		// Menu existed but couldn't be deleted (>48h) → create holder anyway
		// Old menu is garbage, user starts fresh
		if delResult.menuExisted && !delResult.menuDeleted {
			log.Printf("DEBUG: old menu not deleted, creating holder for user %d", userID)
		}

		return bh.sendHolderMessage(ctx, chatID, userID)
	}, th.CommandEqual("start"))
}

// sendHolderMessage sends the welcome holder message with "Start" button.
func (bh *BotHandler) sendHolderMessage(ctx *th.Context, chatID, userID int64) error {
	m := bh.localizer.Get(bh.getLang(userID))
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Start, CallbackData: "holder_start"}},
		},
	}

	sentMsg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		m.HolderText,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard))
	if err != nil {
		log.Printf("ERROR: send holder message: %v", err)
		return err
	}

	if err := bh.sessionStorage.SetHolderMessageID(bh.ctx, userID, sentMsg.MessageID); err != nil {
		log.Printf("ERROR: save holder message id: %v", err)
	}

	return nil
}

func (bh *BotHandler) validateCallback(cb telego.CallbackQuery) (
	chatID int64, messageID int, userID int64, valid bool,
) {
	info := bh.callbackProcessor.Extract(cb)
	if info == nil {
		return 0, 0, 0, false
	}

	bh.rateLimiter.WaitAndGo(bh.ctx, info.UserID)
	if bh.ctx.Err() != nil {
		return 0, 0, 0, false
	}

	return info.ChatID, info.MessageID, info.UserID, true
}

// handleHolderStart processes "Start" button press from holder message.
// Creates new menu and schedules it for cleanup.
func (bh *BotHandler) handleHolderStart(cb telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		return nil
	}
	chatID := msg.Chat.ID
	userID := cb.From.ID
	holderMessageID := msg.MessageID

	if !handlers.AnswerCallbackOrDelete(bh.ctx, bh.bot, cb.ID, chatID, holderMessageID) {
		return nil
	}

	bh.rateLimiter.WaitAndGo(bh.ctx, userID)
	if bh.ctx.Err() != nil {
		return bh.ctx.Err()
	}

	// If menu already exists, try to delete it and remove from queue
	if oldMenuID, err := bh.sessionStorage.GetMenuMessageID(bh.ctx, userID); err == nil && oldMenuID != 0 {
		bh.deleteMessage(chatID, oldMenuID)
		_ = bh.sessionStorage.RemoveFromCleanupQueue(bh.ctx, userID)
		_ = bh.sessionStorage.ClearCleanupRetry(bh.ctx, userID)
	}

	// Note: Holder is NOT deleted — it stays as reference for user

	return bh.sendMenuMessage(chatID, userID)
}

// sendMenuMessage sends the main menu message and schedules cleanup.
func (bh *BotHandler) sendMenuMessage(chatID, userID int64) error {
	m := bh.localizer.Get(bh.getLang(userID))
	sentMsg, err := bh.bot.SendMessage(bh.ctx, tu.Message(
		tu.ID(chatID),
		m.MainMenuText,
	).WithParseMode("Markdown").WithReplyMarkup(bh.getMainMenuInline(m)))
	if err != nil {
		log.Printf("ERROR: send menu message: %v", err)
		return err
	}

	now := time.Now()

	// Save menu message ID
	if err := bh.sessionStorage.SetMenuMessageID(bh.ctx, userID, sentMsg.MessageID); err != nil {
		log.Printf("ERROR: save menu message id: %v", err)
	}

	// Save menu creation time
	if err := bh.sessionStorage.SetMenuCreatedAt(bh.ctx, userID, now); err != nil {
		log.Printf("ERROR: save menu created at: %v", err)
	}

	// Schedule cleanup in 47 hours
	checkAt := now.Add(session.CleanupDelay)
	if err := bh.sessionStorage.AddToCleanupQueue(bh.ctx, userID, checkAt); err != nil {
		log.Printf("ERROR: add to cleanup queue: %v", err)
	}

	return nil
}

func (bh *BotHandler) registerMenuCallbackHandler() {
	// Handle "Start" button from holder message — creates new menu
	bh.handler.HandleCallbackQuery(func(_ *th.Context, cb telego.CallbackQuery) error {
		return bh.handleHolderStart(cb)
	}, th.CallbackDataEqual("holder_start"))

	// Handle back to menu
	bh.handler.HandleCallbackQuery(func(_ *th.Context, cb telego.CallbackQuery) error {
		chatID, messageID, userID, valid := bh.validateCallback(cb)
		if !valid {
			return nil
		}
		bh.showMainMenu(chatID, userID, messageID)
		return nil
	}, th.CallbackDataEqual("menu_back"))
}

func (bh *BotHandler) showMainMenu(chatID, userID int64, messageID int) {
	m := bh.localizer.Get(bh.getLang(userID))

	if _, err := bh.bot.EditMessageText(bh.ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.MainMenuText,
		ParseMode:   "Markdown",
		ReplyMarkup: bh.getMainMenuInline(m),
	}); err != nil {
		if !handlers.HandleEditError(bh.ctx, bh.bot, err, chatID, messageID) {
			log.Printf("ERROR: edit to main menu: %v", err)
		}
	}
}

func (bh *BotHandler) registerTechniqueHandlers() {
	bh.registerBreathingHandler()
	bh.registerGroundingHandler()
	bh.registerGuidedBreathingHandler()
	bh.registerPMRHandler()
	bh.registerThoughtLabelingHandler()
	bh.registerVisualizationHandler()
	bh.registerLangHandler()
}

func (bh *BotHandler) registerBreathingHandler() {
	breathingHandler := handlers.NewBreathingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(breathingHandler.HandleCallback, th.CallbackDataPrefix("breathing_"))
	bh.handler.HandleCallbackQuery(breathingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_breathing"))
}

func (bh *BotHandler) registerGroundingHandler() {
	groundingHandler := handlers.NewGroundingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(groundingHandler.HandleCallback, th.CallbackDataPrefix("grounding_"))
	bh.handler.HandleCallbackQuery(groundingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_grounding"))
}

func (bh *BotHandler) registerGuidedBreathingHandler() {
	guidedBreathingHandler := handlers.NewGuidedBreathingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(guidedBreathingHandler.HandleCallback, th.CallbackDataPrefix("gbreath_"))
	bh.handler.HandleCallbackQuery(guidedBreathingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_guided"))
}

func (bh *BotHandler) registerPMRHandler() {
	pmrHandler := handlers.NewPMRHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(pmrHandler.HandleCallback, th.CallbackDataPrefix("pmr_"))
	bh.handler.HandleCallbackQuery(pmrHandler.HandleMenuSelect, th.CallbackDataEqual("menu_pmr"))
}

func (bh *BotHandler) registerThoughtLabelingHandler() {
	thoughtLabelingHandler := handlers.NewThoughtLabelingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(thoughtLabelingHandler.HandleCallback, th.CallbackDataPrefix("thought_"))
	bh.handler.HandleCallbackQuery(thoughtLabelingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_thought"))
}

func (bh *BotHandler) registerVisualizationHandler() {
	visualizationHandler := handlers.NewVisualizationHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(visualizationHandler.HandleCallback, th.CallbackDataPrefix("visual_"))
	bh.handler.HandleCallbackQuery(visualizationHandler.HandleMenuSelect, th.CallbackDataEqual("menu_visual"))
}

func (bh *BotHandler) registerLangHandler() {
	langHandler := handlers.NewLangHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
		bh.sessionStorage, bh.callbackProcessor,
	)
	bh.handler.HandleCallbackQuery(langHandler.HandleCallback, th.CallbackDataPrefix("lang_"))
	bh.handler.HandleCallbackQuery(langHandler.HandleMenuSelect, th.CallbackDataEqual("menu_lang"))
}

// registerCatchAllHandler handles any unrecognized messages.
func (bh *BotHandler) registerCatchAllHandler() {
	bh.handler.HandleMessage(func(_ *th.Context, message telego.Message) error {
		if message.From == nil {
			return nil
		}
		userID := message.From.ID
		chatID := message.Chat.ID

		// Check if user is in thought labeling input state
		state, _ := bh.sessionStorage.GetState(bh.ctx, userID)
		if state == session.StateThoughtLabelingInput {
			// Get saved message ID to edit
			botMessageID, err := bh.sessionStorage.GetMessageID(bh.ctx, userID)
			if err == nil && botMessageID != 0 {
				// Delete user's message
				bh.deleteMessage(chatID, message.MessageID)

				// Process thought input
				thoughtHandler := handlers.NewThoughtLabelingHandler(
					bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics,
					bh.sessionStorage, bh.sessionManager, bh.callbackProcessor,
				)
				thoughtHandler.ProcessThoughtInput(chatID, userID, botMessageID, message.Text)
				return nil
			}
		}

		// Delete any unrecognized message to keep chat clean
		bh.deleteMessage(chatID, message.MessageID)
		return nil
	}, th.AnyMessage())
}

// deleteMessage silently deletes a message.
func (bh *BotHandler) deleteMessage(chatID int64, messageID int) {
	if err := bh.bot.DeleteMessage(bh.ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	}); err != nil {
		// Ignore errors - message might already be deleted or too old
		log.Printf("DEBUG: could not delete message %d: %v", messageID, err)
	}
}

// Start runs the bot handler in a goroutine.
func (bh *BotHandler) Start() {
	log.Println("Bot started")
	if err := bh.handler.Start(); err != nil {
		log.Printf("ERROR: bot handler start: %v", err)
	}
}

// Stop stops the bot handler gracefully.
func (bh *BotHandler) Stop() {
	if err := bh.handler.Stop(); err != nil {
		log.Printf("ERROR: bot handler stop: %v", err)
	}
	bh.cancelFunc()
	log.Println("Bot stopped")
}

// GetBot returns the underlying telego.Bot instance for use by other services.
func (bh *BotHandler) GetBot() *telego.Bot {
	return bh.bot
}
