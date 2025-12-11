package handlers

import "github.com/mymmrac/telego"

func GetMainMenu() *telego.ReplyKeyboardMarkup {
	return &telego.ReplyKeyboardMarkup{
		Keyboard: [][]telego.KeyboardButton{
			{
				{Text: "🌬️ Дыхание за 2 минуты"},
				{Text: "🌿 Якорение 5-4-3-2-1"},
			},
			{
				{Text: "ℹ️ Информация"},
			},
		},
		ResizeKeyboard: true,
	}
}
