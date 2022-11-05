package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg *telegramBot) onCommandCreate(message *tgbotapi.Message) error {
	switch message.Text {
	case "/start":
		if err := tg.printMessage(message, commandsMessage); err != nil {
			return err
		}
	default:
		if err := tg.printMessage(message, selectPlace); err != nil {
			return err
		}
	}
	return nil
}

func (tg *telegramBot) onCallbackQuery(callback *tgbotapi.CallbackQuery) error {
	if err := tg.callbackRequest(callback); err != nil {
		return err
	}
	return nil
}
