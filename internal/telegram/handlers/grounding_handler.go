package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/db"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
)

type GroundingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	db             db.DBWork
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewGroundingHandler(
	ctx context.Context,
	bot *telego.Bot,
	database db.DBWork,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *GroundingHandler {
	return &GroundingHandler{
		ctx:            ctx,
		bot:            bot,
		db:             database,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

func (h *GroundingHandler) Handle(ctx *th.Context, update telego.Update) error {
	var chatID int64
	var userID int64
	var messageID int

	if update.Message != nil {
		chatID = update.Message.Chat.ID
		userID = update.Message.From.ID
		h.statistics.IncreaseRequestsStatisticForUser(
			userID,
			update.Message.From.Username,
			update.Message.From.IsPremium,
			update.Message.From.IsBot,
		)
		h.startGroundingExercise(h.ctx, chatID, userID)

	} else if update.CallbackQuery != nil {
		msg, ok := update.CallbackQuery.Message.(*telego.Message)
		if !ok || msg == nil {
			log.Printf("ERROR: callback query message is inaccessible")
			return nil
		}

		chatID = msg.Chat.ID
		userID = update.CallbackQuery.From.ID
		messageID = msg.MessageID

		if strings.HasPrefix(update.CallbackQuery.Data, "grounding_step") {
			parts := strings.Split(update.CallbackQuery.Data, "_")
			if len(parts) == 3 {
				stepNum := parts[2]
				h.showGroundingStep(h.ctx, chatID, userID, messageID, stepNum)
			}
		} else if strings.HasPrefix(update.CallbackQuery.Data, "grounding_complete") {
			h.completeGrounding(h.ctx, chatID, userID, messageID)
		} else if update.CallbackQuery.Data == "grounding_cancel" {
			h.cancelGrounding(h.ctx, chatID, userID, messageID)
		}

		_ = ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
	}
	return nil
}

func (h *GroundingHandler) startGroundingExercise(ctx context.Context, chatID, userID int64) {
	_ = h.sessionStorage.SetState(ctx, userID, session.StateGroundingStep1)

	intro := `🌿 *Якорение 5-4-3-2-1*

Техника для возвращения в настоящий момент.

*Суть:* Назовите вслух или про себя:

• 5 вещей, которые видите 👁️
• 4 вещи, которые ощущаете 🤚
• 3 звука, которые слышите 👂
• 2 запаха, которые чувствуете 👃
• 1 вкус во рту 👅

Готовы начать?`

	inlineKeyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "grounding_step_1"},
				{Text: "❌ Отмена", CallbackData: "grounding_cancel"},
			},
		},
	}

	removeKeyboard := &telego.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}

	_, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		intro,
	).WithParseMode("Markdown").WithReplyMarkup(inlineKeyboard))
	if err != nil {
		log.Printf("ERROR: send grounding intro: %v", err)
		return
	}

	_, _ = h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"_Используйте кнопки выше_",
	).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard))
}

func (h *GroundingHandler) showGroundingStep(ctx context.Context, chatID, userID int64, messageID int, stepNum string) {
	var state session.State
	switch stepNum {
	case "1":
		state = session.StateGroundingStep1
	case "2":
		state = session.StateGroundingStep2
	case "3":
		state = session.StateGroundingStep3
	case "4":
		state = session.StateGroundingStep4
	case "5":
		state = session.StateGroundingStep5
	}
	_ = h.sessionStorage.SetState(ctx, userID, state)

	steps := techniques.GetGroundingSteps()
	var step techniques.GroundingStep
	var nextStep string

	switch stepNum {
	case "1":
		step = steps[0]
		nextStep = "2"
	case "2":
		step = steps[1]
		nextStep = "3"
	case "3":
		step = steps[2]
		nextStep = "4"
	case "4":
		step = steps[3]
		nextStep = "5"
	case "5":
		step = steps[4]
		nextStep = "complete"
	}

	text := fmt.Sprintf("%s *%s*\n\n%s", step.Emoji, step.Title, step.Description)

	var keyboard *telego.InlineKeyboardMarkup
	if nextStep == "complete" {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: "✅ Завершить", CallbackData: "grounding_complete"},
				},
			},
		}
	} else {
		keyboard = &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					{Text: "➡️ Далее", CallbackData: "grounding_step_" + nextStep},
				},
			},
		}
	}

	_, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Printf("ERROR: edit grounding message: %v", err)
	}
}

func (h *GroundingHandler) completeGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	_ = h.sessionStorage.ClearState(ctx, userID)

	text := `✨ *Отлично!*

Вы завершили технику якорения 5-4-3-2-1.

Эта практика помогает выйти из тревожных мыслей и вернуться в настоящий момент.`

	_, _ = h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	})

	_, _ = h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu()))
}

func (h *GroundingHandler) cancelGrounding(ctx context.Context, chatID, userID int64, messageID int) {
	_ = h.sessionStorage.ClearState(ctx, userID)

	_, _ = h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      "❌ Упражнение отменено.",
	})

	_, _ = h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu()))
}
