package handlers

import (
	"github.com/mymmrac/telego"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
)

// MainMenuText is the text shown in main menu
var MainMenuText = messages.MainMenuText

// EmptyInlineKeyboard returns empty inline keyboard to remove buttons when editing message.
func EmptyInlineKeyboard() *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{}}
}

// GetMainMenuInline returns inline keyboard for technique selection.
func GetMainMenuInline() *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "🌬️ Дыхание 2 мин", CallbackData: "menu_breathing"},
				{Text: "🌿 Якорение", CallbackData: "menu_grounding"},
			},
			{
				{Text: "🧘 Управляемое", CallbackData: "menu_guided"},
				{Text: "💪 Мышечная", CallbackData: "menu_pmr"},
			},
			{
				{Text: "🏷️ Мысли", CallbackData: "menu_thought"},
				{Text: "🌅 Визуализация", CallbackData: "menu_visual"},
			},
			{
				{Text: "ℹ️ Информация", CallbackData: "menu_info"},
			},
		},
	}
}
