package helpers

import tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// CreateKeyboard
// create keybowrds of two rows of any map[string]string input
func CreateKeyboard(data map[string]string) tgb.InlineKeyboardMarkup {
	// hardcoded models
	keyboard := tgb.NewInlineKeyboardMarkup()
	//	subbuttons := []tgbot.InlineKeyboardButton{}
	rows := tgb.NewInlineKeyboardRow()
	counter := 0
	for key, val := range data {

		if counter != 0 && counter%3 == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rows)
			rows = tgb.NewInlineKeyboardRow()
		}
		rows = append(rows, tgb.NewInlineKeyboardButtonData(key, val))
		if counter >= len(data)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rows)
		}
		counter++
	}
	return keyboard
}
