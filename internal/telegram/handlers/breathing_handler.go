package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/db"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
)

type BreathingHandler struct {
	ctx         context.Context
	bot         *telego.Bot
	db          db.DBWork
	rateLimiter rate_limiter.Limiter
	statistics  statistic.Stats
}

func NewBreathingHandler(
	ctx context.Context,
	bot *telego.Bot,
	database db.DBWork,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
) *BreathingHandler {
	return &BreathingHandler{
		ctx:         ctx,
		bot:         bot,
		db:          database,
		rateLimiter: rateLimiter,
		statistics:  statistics,
	}
}

func (h *BreathingHandler) Handle(ctx *th.Context, update telego.Update) error {
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
		h.startBreathingExercise(h.ctx, chatID, userID)

	} else if update.CallbackQuery != nil {
		msg, ok := update.CallbackQuery.Message.(*telego.Message)
		if !ok || msg == nil {
			log.Printf("ERROR: callback query message is inaccessible")
			return nil
		}

		chatID = msg.Chat.ID
		userID = update.CallbackQuery.From.ID
		messageID = msg.MessageID

		if strings.HasPrefix(update.CallbackQuery.Data, "breathing_complete") {
			h.completeBreathing(h.ctx, chatID, userID, messageID)
		} else if strings.HasPrefix(update.CallbackQuery.Data, "breathing_cancel") {
			h.cancelBreathing(h.ctx, chatID, messageID)
		} else if update.CallbackQuery.Data == "breathing_start" {
			go h.runBreathingCycle(chatID, userID, messageID)
		}

		_ = ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
	}
	return nil
}

func (h *BreathingHandler) startBreathingExercise(ctx context.Context, chatID, userID int64) {
	intro := `🌬️ *Дыхание за 2 минуты*

Простое упражнение для успокоения нервной системы.

*Инструкция:*
1️⃣ Сядьте удобно, закройте глаза
2️⃣ Вдох через нос — 4 секунды
3️⃣ Задержка дыхания — 4 секунды
4️⃣ Выдох через рот — 6 секунд
5️⃣ Повторите 8 циклов (~2 минуты)

Я буду напоминать каждый этап. Готовы начать?`

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "breathing_start"},
			},
			{
				{Text: "❌ Отмена", CallbackData: "breathing_cancel"},
			},
		},
	}

	_, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		intro,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard))
	if err != nil {
		log.Printf("ERROR: send breathing intro: %v", err)
	}
}

func (h *BreathingHandler) runBreathingCycle(chatID int64, userID int64, messageID int) {
	ctx := h.ctx
	cycles := techniques.GetBreathingCycles()

	for i, cycle := range cycles {
		select {
		case <-ctx.Done():
			return
		default:
			text := fmt.Sprintf("🌬️ *Цикл %d из 8*\n\n%s", i+1, cycle.Instruction)

			_, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
				ChatID:    tu.ID(chatID),
				MessageID: messageID,
				Text:      text,
				ParseMode: "Markdown",
			})
			if err != nil {
				log.Printf("ERROR: edit breathing message: %v", err)
			}

			time.Sleep(cycle.Duration)
		}
	}

	completionText := `✅ *Отлично!*

Вы завершили дыхательное упражнение.
Как вы себя чувствуете?`

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "✅ Лучше", CallbackData: "breathing_complete"},
				{Text: "🔄 Повторить", CallbackData: "breathing_start"},
			},
		},
	}

	_, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        completionText,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Printf("ERROR: send completion message: %v", err)
	}
}

func (h *BreathingHandler) completeBreathing(ctx context.Context, chatID int64, userID int64, messageID int) {
	text := `✨ *Спасибо за практику!*

Регулярные упражнения помогают снизить уровень тревожности.`

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

func (h *BreathingHandler) cancelBreathing(ctx context.Context, chatID int64, messageID int) {
	_, _ = h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      "❌ Упражнение отменено.",
	})
}
