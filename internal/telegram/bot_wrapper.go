// Package telegram implements the Telegram bot handlers and routing.
package telegram

import (
	"context"
	"log"

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
	ctx            context.Context
	cancelFunc     context.CancelFunc
	bot            *telego.Bot
	handler        *th.BotHandler
	localizer      *localization.Localizer
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
	sessionManager *session.SessionManager
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

	bh := &BotHandler{
		ctx:            ctx,
		cancelFunc:     cancel,
		bot:            bot,
		handler:        botHandler,
		localizer:      localization.NewLocalizer(),
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
		sessionManager: session.NewSessionManager(),
	}

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
				{Text: m.MenuInfo, CallbackData: "menu_info"},
			},
			{
				{Text: m.MenuLang, CallbackData: "menu_lang"},
			},
		},
	}
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

		// Cancel any running session (stops active goroutines)
		bh.sessionManager.CancelSession(userID)

		// Delete old bot message if exists
		if oldMsgID, err := bh.sessionStorage.GetMessageID(bh.ctx, userID); err == nil && oldMsgID != 0 {
			bh.deleteMessage(chatID, oldMsgID)
		}

		// Clear session state
		if err := bh.sessionStorage.ClearState(bh.ctx, userID); err != nil {
			log.Printf("ERROR: clear state on start: %v", err)
		}

		bh.statistics.IncreaseRequestsStatisticForUser(
			userID,
			message.From.Username,
			message.From.IsPremium,
			message.From.IsBot,
		)

		m := bh.localizer.Get(bh.getLang(userID))
		welcomeText := "👋 *" + message.From.FirstName + "*!\n\n" + m.MainMenuText

		sentMsg, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(chatID),
			welcomeText,
		).WithParseMode("Markdown").WithReplyMarkup(bh.getMainMenuInline(m)))
		if err != nil {
			log.Printf("ERROR: send start message: %v", err)
			return err
		}

		// Save new message ID for future deletion
		if err := bh.sessionStorage.SetMessageID(bh.ctx, userID, sentMsg.MessageID); err != nil {
			log.Printf("ERROR: save message id: %v", err)
		}

		return nil
	}, th.CommandEqual("start"))
}

func (bh *BotHandler) registerMenuCallbackHandler() {
	// Handle info
	bh.handler.HandleCallbackQuery(func(_ *th.Context, cb telego.CallbackQuery) error {
		msg, ok := cb.Message.(*telego.Message)
		if !ok || msg == nil {
			return nil
		}
		chatID := msg.Chat.ID
		messageID := msg.MessageID

		// Answer callback FIRST; if too old - delete message and stop
		if !handlers.AnswerCallbackOrDelete(bh.ctx, bh.bot, cb.ID, chatID, messageID) {
			return nil
		}

		bh.rateLimiter.WaitAndGo(bh.ctx, cb.From.ID)
		if bh.ctx.Err() != nil {
			return bh.ctx.Err()
		}

		bh.showInfo(chatID, cb.From.ID, messageID)
		return nil
	}, th.CallbackDataEqual("menu_info"))

	// Handle back to menu
	bh.handler.HandleCallbackQuery(func(_ *th.Context, cb telego.CallbackQuery) error {
		msg, ok := cb.Message.(*telego.Message)
		if !ok || msg == nil {
			return nil
		}
		chatID := msg.Chat.ID
		messageID := msg.MessageID

		// Answer callback FIRST; if too old - delete message and stop
		if !handlers.AnswerCallbackOrDelete(bh.ctx, bh.bot, cb.ID, chatID, messageID) {
			return nil
		}

		bh.rateLimiter.WaitAndGo(bh.ctx, cb.From.ID)
		if bh.ctx.Err() != nil {
			return bh.ctx.Err()
		}

		bh.showMainMenu(chatID, cb.From.ID, messageID)
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

func (bh *BotHandler) showInfo(chatID, userID int64, messageID int) {
	m := bh.localizer.Get(bh.getLang(userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.BackToMenu, CallbackData: "menu_back"}},
		},
	}

	if _, err := bh.bot.EditMessageText(bh.ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        m.InfoText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit to info: %v", err)
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
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
	)
	bh.handler.HandleCallbackQuery(breathingHandler.HandleCallback, th.CallbackDataPrefix("breathing_"))
	bh.handler.HandleCallbackQuery(breathingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_breathing"))
}

func (bh *BotHandler) registerGroundingHandler() {
	groundingHandler := handlers.NewGroundingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
	)
	bh.handler.HandleCallbackQuery(groundingHandler.HandleCallback, th.CallbackDataPrefix("grounding_"))
	bh.handler.HandleCallbackQuery(groundingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_grounding"))
}

func (bh *BotHandler) registerGuidedBreathingHandler() {
	guidedBreathingHandler := handlers.NewGuidedBreathingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
	)
	bh.handler.HandleCallbackQuery(guidedBreathingHandler.HandleCallback, th.CallbackDataPrefix("gbreath_"))
	bh.handler.HandleCallbackQuery(guidedBreathingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_guided"))
}

func (bh *BotHandler) registerPMRHandler() {
	pmrHandler := handlers.NewPMRHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
	)
	bh.handler.HandleCallbackQuery(pmrHandler.HandleCallback, th.CallbackDataPrefix("pmr_"))
	bh.handler.HandleCallbackQuery(pmrHandler.HandleMenuSelect, th.CallbackDataEqual("menu_pmr"))
}

func (bh *BotHandler) registerThoughtLabelingHandler() {
	thoughtLabelingHandler := handlers.NewThoughtLabelingHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
	)
	bh.handler.HandleCallbackQuery(thoughtLabelingHandler.HandleCallback, th.CallbackDataPrefix("thought_"))
	bh.handler.HandleCallbackQuery(thoughtLabelingHandler.HandleMenuSelect, th.CallbackDataEqual("menu_thought"))
}

func (bh *BotHandler) registerVisualizationHandler() {
	visualizationHandler := handlers.NewVisualizationHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
	)
	bh.handler.HandleCallbackQuery(visualizationHandler.HandleCallback, th.CallbackDataPrefix("visual_"))
	bh.handler.HandleCallbackQuery(visualizationHandler.HandleMenuSelect, th.CallbackDataEqual("menu_visual"))
}

func (bh *BotHandler) registerLangHandler() {
	langHandler := handlers.NewLangHandler(
		bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage,
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
					bh.ctx, bh.bot, bh.localizer, bh.rateLimiter, bh.statistics, bh.sessionStorage, bh.sessionManager,
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
