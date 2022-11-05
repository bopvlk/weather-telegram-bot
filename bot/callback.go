package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg *telegramBot) callbackRequest(callback *tgbotapi.CallbackQuery) error {
	clbck := tgbotapi.NewCallback(callback.ID, callback.Data)
	if _, err := tg.bot.Request(clbck); err != nil {
		return fmt.Errorf("send request failed: %v", err)
	}
	if clbck.Text[0] == 'l' && clbck.Text[1] == 'a' && clbck.Text[2] == 't' {
		forecast, err := tg.forecastRequest(clbck.Text)
		if err != nil {
			return err
		}
		if err := tg.printMessage(callback.Message, forecast); err != nil {
			return err
		}
	}
	return nil
}

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
