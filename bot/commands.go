package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandsMessage = "Hello I am an Forecast bot. Enter the name of the city and you will get forecast."
	uknownMessage   = "Uknown Command!!  Please write /start"
	selectPlace     = "Select a right place"
)

func (tg *telegramBot) printMessage(message *tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if text == selectPlace {
		tg.setGeolocationRequest(message.Text)
		msg.ReplyMarkup = tg.citySelection()
	}
	if err := tg.sendMessage(msg); err != nil {
		return err
	}
	return nil
}

func (tg *telegramBot) sendMessage(msg tgbotapi.MessageConfig) error {
	if _, err := tg.bot.Send(msg); err != nil {
		return fmt.Errorf("send message failed: %v", err)
	}
	return nil
}
