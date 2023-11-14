package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	keyStartAuthorization = "Authorize for faster work."
	keyStartCheckOut      = "Check out the weather forecast right now"
	keyJustLoggedYES      = "Yes!"
	keyJustLoggedSchedule = "I wanna make forecast schedule"
)

var (
	keyboardStarted = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(keyStartCheckOut, keyStartCheckOut)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(keyStartAuthorization, keyStartAuthorization)))
	keyboardJustLogged = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(keyJustLoggedYES, keyJustLoggedYES)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(keyJustLoggedSchedule, keyJustLoggedSchedule)),
	)
)

func (tg *telegramBot) citySelection() tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, g := range tg.forecast.g {
		key := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Name of city: %s, region: %s, country: %s", g.Name, g.Region, g.Country),
			fmt.Sprintf("lat=%.2f&lon=%.2f", g.Latitude, g.Longitude)))
		keyboard = append(keyboard, key)
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
