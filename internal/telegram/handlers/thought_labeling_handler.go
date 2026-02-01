package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/rate_limiter"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
)

type ThoughtLabelingHandler struct {
	ctx            context.Context
	bot            *telego.Bot
	rateLimiter    rate_limiter.Limiter
	statistics     statistic.Stats
	sessionStorage session.Storage
}

func NewThoughtLabelingHandler(
	ctx context.Context,
	bot *telego.Bot,
	rateLimiter rate_limiter.Limiter,
	statistics statistic.Stats,
	sessionStorage session.Storage,
) *ThoughtLabelingHandler {
	return &ThoughtLabelingHandler{
		ctx:            ctx,
		bot:            bot,
		rateLimiter:    rateLimiter,
		statistics:     statistics,
		sessionStorage: sessionStorage,
	}
}

func (h *ThoughtLabelingHandler) Handle(ctx *th.Context, update telego.Update) error {
	if update.Message != nil {
		return h.handleMessage(update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}
	return nil
}

func (h *ThoughtLabelingHandler) handleMessage(msg *telego.Message) error {
	userID := msg.From.ID
	h.rateLimiter.WaitAndGo(h.ctx, userID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}
	h.statistics.IncreaseRequestsStatisticForUser(
		userID,
		msg.From.Username,
		msg.From.IsPremium,
		msg.From.IsBot,
	)

	state, _ := h.sessionStorage.GetState(h.ctx, userID)
	if state == session.StateThoughtLabelingInput {
		h.showCategories(h.ctx, msg.Chat.ID, userID, msg.Text)
		return nil
	}

	h.showIntro(h.ctx, msg.Chat.ID, userID)
	return nil
}

func (h *ThoughtLabelingHandler) handleCallback(ctx *th.Context, cb *telego.CallbackQuery) error {
	msg, ok := cb.Message.(*telego.Message)
	if !ok || msg == nil {
		log.Printf("ERROR: callback query message is inaccessible")
		return nil
	}
	chatID := msg.Chat.ID
	userID := cb.From.ID
	messageID := msg.MessageID
	h.rateLimiter.WaitAndGo(h.ctx, userID)
	if h.ctx.Err() != nil {
		return h.ctx.Err()
	}

	h.processCallback(chatID, userID, messageID, cb.Data)

	if err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cb.ID,
	}); err != nil {
		log.Printf("ERROR: answer callback query: %v", err)
		return err
	}
	return nil
}

func (h *ThoughtLabelingHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch {
	case data == "thought_start":
		h.promptForThought(h.ctx, chatID, userID, messageID)
	case strings.HasPrefix(data, "thought_cat_"):
		categoryID := strings.TrimPrefix(data, "thought_cat_")
		h.showCategoryInfo(h.ctx, chatID, userID, messageID, categoryID)
	case data == "thought_another":
		h.promptForThought(h.ctx, chatID, userID, messageID)
	case data == "thought_complete":
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case data == "thought_cancel":
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *ThoughtLabelingHandler) showIntro(ctx context.Context, chatID, userID int64) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingActive); err != nil {
		log.Printf("ERROR: set state thought labeling active: %v", err)
	}

	text := `🏷️ *Маркировка мыслей*

Техника когнитивной терапии для работы с тревожными мыслями.

*Как это работает:*
1. Вы описываете тревожную мысль
2. Определяете её тип (категорию)
3. Осознание типа мысли снижает её влияние

*Почему это помогает:*
Когда мы называем мысль, мы отделяем себя от неё. Это не "я тревожусь", а "это мысль типа беспокойство".

Готовы попробовать?`

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "▶️ Начать", CallbackData: "thought_start"},
				{Text: "❌ Отмена", CallbackData: "thought_cancel"},
			},
		},
	}

	removeKeyboard := &telego.ReplyKeyboardRemove{RemoveKeyboard: true}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		text,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard)); err != nil {
		log.Printf("ERROR: send thought labeling intro: %v", err)
		return
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"_Используйте кнопки выше_",
	).WithParseMode("Markdown").WithReplyMarkup(removeKeyboard)); err != nil {
		log.Printf("ERROR: send thought labeling use buttons: %v", err)
	}
}

func (h *ThoughtLabelingHandler) promptForThought(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingInput); err != nil {
		log.Printf("ERROR: set state thought labeling input: %v", err)
	}

	text := `💭 *Опишите тревожную мысль*

Напишите мысль, которая вас беспокоит.

_Например: "Я боюсь, что не справлюсь с работой"_`

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: "❌ Отмена", CallbackData: "thought_cancel"}},
		},
	}

	if messageID == 0 {
		if _, err := h.bot.SendMessage(ctx, tu.Message(
			tu.ID(chatID),
			text,
		).WithParseMode("Markdown").WithReplyMarkup(keyboard)); err != nil {
			log.Printf("ERROR: send thought prompt: %v", err)
		}
	} else {
		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        text,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil {
			log.Printf("ERROR: edit thought prompt: %v", err)
		}
	}
}

func (h *ThoughtLabelingHandler) showCategories(
	ctx context.Context, chatID, userID int64, thought string,
) {
	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingCategory); err != nil {
		log.Printf("ERROR: set state thought labeling category: %v", err)
	}

	categories := techniques.GetThoughtCategories()

	text := fmt.Sprintf("💭 *Ваша мысль:*\n_%s_\n\n🏷️ *К какому типу относится эта мысль?*", thought)

	buttons := h.buildCategoryButtons(categories)
	keyboard := &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		text,
	).WithParseMode("Markdown").WithReplyMarkup(keyboard)); err != nil {
		log.Printf("ERROR: send thought categories: %v", err)
	}
}

func (h *ThoughtLabelingHandler) buildCategoryButtons(
	categories []techniques.ThoughtCategory,
) [][]telego.InlineKeyboardButton {
	var buttons [][]telego.InlineKeyboardButton
	row := []telego.InlineKeyboardButton{}
	for i, cat := range categories {
		row = append(row, telego.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s %s", cat.Emoji, cat.Name),
			CallbackData: fmt.Sprintf("thought_cat_%s", cat.ID),
		})
		if len(row) == 2 || i == len(categories)-1 {
			buttons = append(buttons, row)
			row = []telego.InlineKeyboardButton{}
		}
	}
	buttons = append(buttons, []telego.InlineKeyboardButton{
		{Text: "❌ Отмена", CallbackData: "thought_cancel"},
	})
	return buttons
}

func (h *ThoughtLabelingHandler) showCategoryInfo(
	ctx context.Context, chatID, userID int64, messageID int, categoryID string,
) {
	categories := techniques.GetThoughtCategories()
	var category *techniques.ThoughtCategory
	for _, cat := range categories {
		if cat.ID == categoryID {
			category = &cat
			break
		}
	}
	if category == nil {
		return
	}

	if err := h.sessionStorage.SetState(ctx, userID, session.StateThoughtLabelingActive); err != nil {
		log.Printf("ERROR: set state thought labeling active: %v", err)
	}

	text := fmt.Sprintf(`%s *%s*

_%s_

*Пример:* %s

---

✅ Вы успешно промаркировали мысль!

Осознание типа мысли — первый шаг к снижению её влияния. Теперь вы видите эту мысль со стороны.

Хотите промаркировать ещё одну мысль?`, category.Emoji, category.Name, category.Description, category.Example)

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "💭 Ещё мысль", CallbackData: "thought_another"},
				{Text: "✅ Готово", CallbackData: "thought_complete"},
			},
		},
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(chatID),
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "Markdown",
		ReplyMarkup: keyboard,
	}); err != nil {
		log.Printf("ERROR: edit thought category info: %v", err)
	}
}

func (h *ThoughtLabelingHandler) completeExercise(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	text := `✨ *Спасибо за практику!*

Маркировка мыслей — мощная техника когнитивной терапии. Регулярная практика помогает:

• Снижать влияние тревожных мыслей
• Развивать осознанность
• Отделять себя от негативных паттернов мышления`

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}); err != nil {
		log.Printf("ERROR: edit thought complete: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Выберите другую технику:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send thought menu: %v", err)
	}
}

func (h *ThoughtLabelingHandler) cancelExercise(
	ctx context.Context, chatID, userID int64, messageID int,
) {
	if err := h.sessionStorage.ClearState(ctx, userID); err != nil {
		log.Printf("ERROR: clear state: %v", err)
	}

	if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      "❌ Упражнение отменено.",
	}); err != nil {
		log.Printf("ERROR: edit thought cancel: %v", err)
	}

	if _, err := h.bot.SendMessage(ctx, tu.Message(
		tu.ID(chatID),
		"Возврат в главное меню:",
	).WithReplyMarkup(GetMainMenu())); err != nil {
		log.Printf("ERROR: send thought menu: %v", err)
	}
}
