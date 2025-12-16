package telegram

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/db"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/handlers"
)

type BotHandler struct {
	ctx            context.Context
	cancelFunc     context.CancelFunc
	bot            *telego.Bot
	handler        *th.BotHandler
	db             db.DBWork
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func MustNewBotHandler(
	ctx context.Context,
	token string,
	database db.DBWork,
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
		db:             database,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}

	bh.registerHandlers()

	return bh
}

func (bh *BotHandler) registerHandlers() {
	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		bh.statistics.IncreaseRequestsStatisticForUser(
			message.From.ID,
			message.From.Username,
			message.From.IsPremium,
			message.From.IsBot,
		)

		welcomeText := `👋 Привет, *` + message.From.FirstName + `*!

Я — AnxietyHelp бот, помогу справиться с тревожностью.

🌬️ *Дыхание за 2 минуты* — быстрое успокоение
🌿 *Якорение 5-4-3-2-1* — заземление в настоящем

Выберите технику из меню ниже:`

		_, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			welcomeText,
		).WithParseMode("Markdown").WithReplyMarkup(handlers.GetMainMenu()))
		if err != nil {
			log.Printf("ERROR: send start message: %v", err)
		}
		return nil
	}, th.CommandEqual("start"))

	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		bh.statistics.IncreaseRequestsStatisticForUser(
			message.From.ID,
			message.From.Username,
			message.From.IsPremium,
			message.From.IsBot,
		)

		menuText := `🏠 *Главное меню*

Выберите технику для работы с тревожностью:`

		_, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			menuText,
		).WithParseMode("Markdown").WithReplyMarkup(handlers.GetMainMenu()))
		if err != nil {
			log.Printf("ERROR: send menu: %v", err)
		}
		return nil
	}, th.TextEqual("🏠 Главное меню"))

	breathingHandler := handlers.NewBreathingHandler(bh.ctx, bh.bot, bh.db, bh.rateLimiter, bh.statistics, bh.sessionStorage)
	bh.handler.Handle(breathingHandler.Handle, th.Or(
		th.TextEqual("🌬️ Дыхание за 2 минуты"),
		th.CallbackDataPrefix("breathing_"),
	))

	groundingHandler := handlers.NewGroundingHandler(bh.ctx, bh.bot, bh.db, bh.rateLimiter, bh.statistics, bh.sessionStorage)
	bh.handler.Handle(groundingHandler.Handle, th.Or(
		th.TextEqual("🌿 Якорение 5-4-3-2-1"),
		th.CallbackDataPrefix("grounding_"),
	))

	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		infoText := `ℹ️ *О боте AnxietyHelp*

Этот бот предоставляет простые техники для снижения тревожности:

🌬️ *Дыхание за 2 минуты* — успокаивает нервную систему
🌿 *Якорение 5-4-3-2-1* — возвращает в настоящий момент

💡 *Совет:* Практикуйте регулярно для лучшего эффекта.

⚠️ При постоянной тревоге обратитесь к специалисту.`

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(message.Chat.ID), infoText).WithParseMode("Markdown"))
		return nil
	}, th.TextEqual("ℹ️ Информация"))
}

func (bh *BotHandler) Start() {
	log.Println("Bot started")
	bh.handler.Start()
}

func (bh *BotHandler) Stop() {
	bh.handler.Stop()
	bh.cancelFunc()
	log.Println("Bot stopped")
}
