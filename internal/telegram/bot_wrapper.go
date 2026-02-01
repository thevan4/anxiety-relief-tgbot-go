package telegram

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
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
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

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
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}

	bh.registerHandlers()

	return bh
}

func (bh *BotHandler) registerHandlers() {
	bh.registerStartAndMenu()
	bh.registerTechniqueHandlers()
	bh.registerInfoHandler()
}

func (bh *BotHandler) registerStartAndMenu() {
	bh.registerStartHandler()
	bh.registerMenuHandler()
}

func (bh *BotHandler) registerStartHandler() {
	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if message.From == nil {
			log.Printf("ERROR: start message has no From")
			return nil
		}
		bh.rateLimiter.WaitAndGo(bh.ctx, message.From.ID)
		if bh.ctx.Err() != nil {
			return bh.ctx.Err()
		}
		bh.statistics.IncreaseRequestsStatisticForUser(
			message.From.ID,
			message.From.Username,
			message.From.IsPremium,
			message.From.IsBot,
		)
		welcomeText := `👋 Привет, *` + message.From.FirstName + `*!

Я — AnxietyHelp бот, помогу справиться с тревожностью.

*Быстрые техники:*
🌬️ Дыхание за 2 минуты
🌿 Якорение 5-4-3-2-1

*Продвинутые техники:*
🧘 Управляемое дыхание
💪 Мышечная релаксация
🏷️ Маркировка мыслей
🌅 Визуализация

Выберите технику из меню ниже:`
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			welcomeText,
		).WithParseMode("Markdown").WithReplyMarkup(handlers.GetMainMenu()))
		if err != nil {
			log.Printf("ERROR: send start message: %v", err)
			return err
		}
		return nil
	}, th.CommandEqual("start"))
}

func (bh *BotHandler) registerMenuHandler() {
	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if message.From == nil {
			log.Printf("ERROR: menu message has no From")
			return nil
		}
		bh.rateLimiter.WaitAndGo(bh.ctx, message.From.ID)
		if bh.ctx.Err() != nil {
			return bh.ctx.Err()
		}
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
			return err
		}
		return nil
	}, th.TextEqual("🏠 Главное меню"))
}

func (bh *BotHandler) registerTechniqueHandlers() {
	bh.registerBreathingHandler()
	bh.registerGroundingHandler()
	bh.registerGuidedBreathingHandler()
	bh.registerPMRHandler()
	bh.registerThoughtLabelingHandler()
	bh.registerVisualizationHandler()
}

func (bh *BotHandler) registerBreathingHandler() {
	breathingHandler := handlers.NewBreathingHandler(
		bh.ctx, bh.bot, bh.rateLimiter, bh.statistics, bh.sessionStorage,
	)
	bh.handler.Handle(breathingHandler.Handle, th.Or(
		th.TextEqual("🌬️ Дыхание за 2 минуты"),
		th.CallbackDataPrefix("breathing_"),
	))
}

func (bh *BotHandler) registerGroundingHandler() {
	groundingHandler := handlers.NewGroundingHandler(
		bh.ctx, bh.bot, bh.rateLimiter, bh.statistics, bh.sessionStorage,
	)
	bh.handler.Handle(groundingHandler.Handle, th.Or(
		th.TextEqual("🌿 Якорение 5-4-3-2-1"),
		th.CallbackDataPrefix("grounding_"),
	))
}

func (bh *BotHandler) registerGuidedBreathingHandler() {
	guidedBreathingHandler := handlers.NewGuidedBreathingHandler(
		bh.ctx, bh.bot, bh.rateLimiter, bh.statistics, bh.sessionStorage,
	)
	bh.handler.Handle(guidedBreathingHandler.Handle, th.Or(
		th.TextEqual("🧘 Управляемое дыхание"),
		th.CallbackDataPrefix("gbreath_"),
	))
}

func (bh *BotHandler) registerPMRHandler() {
	pmrHandler := handlers.NewPMRHandler(
		bh.ctx, bh.bot, bh.rateLimiter, bh.statistics, bh.sessionStorage,
	)
	bh.handler.Handle(pmrHandler.Handle, th.Or(
		th.TextEqual("💪 Мышечная релаксация"),
		th.CallbackDataPrefix("pmr_"),
	))
}

func (bh *BotHandler) registerThoughtLabelingHandler() {
	thoughtLabelingHandler := handlers.NewThoughtLabelingHandler(
		bh.ctx, bh.bot, bh.rateLimiter, bh.statistics, bh.sessionStorage,
	)
	bh.handler.Handle(thoughtLabelingHandler.Handle, th.Or(
		th.TextEqual("🏷️ Маркировка мыслей"),
		th.CallbackDataPrefix("thought_"),
	))
}

func (bh *BotHandler) registerVisualizationHandler() {
	visualizationHandler := handlers.NewVisualizationHandler(
		bh.ctx, bh.bot, bh.rateLimiter, bh.statistics, bh.sessionStorage,
	)
	bh.handler.Handle(visualizationHandler.Handle, th.Or(
		th.TextEqual("🌅 Визуализация"),
		th.CallbackDataPrefix("visual_"),
	))
}

func (bh *BotHandler) registerInfoHandler() {
	bh.handler.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		if message.From == nil {
			log.Printf("ERROR: info message has no From")
			return nil
		}
		bh.rateLimiter.WaitAndGo(bh.ctx, message.From.ID)
		if bh.ctx.Err() != nil {
			return bh.ctx.Err()
		}
		infoText := `ℹ️ *О боте AnxietyHelp*

Этот бот предоставляет техники для снижения тревожности:

*Быстрые техники (2-5 мин):*
🌬️ *Дыхание* — успокаивает нервную систему
🌿 *Якорение* — возвращает в настоящий момент

*Продвинутые техники (5-15 мин):*
🧘 *Управляемое дыхание* — разные паттерны дыхания
💪 *Мышечная релаксация* — снятие телесного напряжения
🏷️ *Маркировка мыслей* — работа с тревожными мыслями
🌅 *Визуализация* — расслабление через воображение

💡 *Совет:* Практикуйте регулярно для лучшего эффекта.

⚠️ При постоянной тревоге обратитесь к специалисту.`
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			infoText,
		).WithParseMode("Markdown"))
		if err != nil {
			log.Printf("ERROR: send info message: %v", err)
			return err
		}
		return nil
	}, th.TextEqual("ℹ️ Информация"))
}

// Start runs the bot handler in a goroutine.
func (bh *BotHandler) Start() {
	log.Println("Bot started")
	if err := bh.handler.Start(); err != nil {
		log.Printf("ERROR: bot handler start: %v", err)
	}
}

func (bh *BotHandler) Stop() {
	if err := bh.handler.Stop(); err != nil {
		log.Printf("ERROR: bot handler stop: %v", err)
	}
	bh.cancelFunc()
	log.Println("Bot stopped")
}
