package bot

import (
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	fromBotStartIfLogged   = "Hello I am a Forecast bot. U are logged.\n Do u wanna know forecast?"
	fromBotStartIfNoLogged = "Hello I am a Forecast bot. U aren`t logged."
	fromBotUknownMessage   = "Uknown Command!!  Please write /start"
	fromBotSelectPlace     = "Select a right place"
	fromBotWriteCity       = "Write the name of the city and you will get forecast."
	fromBotNotFoundCity    = "City does not found!\nWrite the name of the city again."
	fromBotPasword         = "Write a password to create an account"
	fromBotWitchCity       = "At the Which city do you want to know the weather forecast?"
	fromBotULogged         = "U are logged in! Do you wanna know forecast?"
	fromBotWhatTime        = "Please write the time, when you want message with forecast. \nTime must be in 24h format, like \"15:25\""
	fromBotWrongFormatTime = "Wrong format time. Please write the time again"
	fromBotScheduleName    = "Plese create name to this schedule"
)

var (
	toDBPasswordHash = ""
	toDBEventTime    = ""
)

func (tg *telegramBot) printMessage(message *tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	switch msg.Text {
	case fromBotStartIfNoLogged:
		msg.ReplyMarkup = keyboardStarted
	case fromBotStartIfLogged:
		msg.ReplyMarkup = keyboardJustLogged
	case fromBotSelectPlace:
		if err := tg.setGeolocationRequest(message.Text); err != nil {
			return err
		}
		if len(tg.forecast.g) == 0 {
			msg.Text = fromBotNotFoundCity
			tg.pageMarker[message.From.ID] = models.Pages{MarkerFindCity: true}
		} else {
			msg.ReplyMarkup = tg.citySelection()
		}
	case fromBotWriteCity:
		tg.pageMarker[message.From.ID] = models.Pages{MarkerFindCity: true}
	case fromBotPasword:
		tg.pageMarker[message.From.ID] = models.Pages{MarkerWritePassword: true}
	case fromBotWitchCity:
		tg.pageMarker[message.From.ID] = models.Pages{MarkerSaveCityMarker: true, MarkerFindCity: true}
	case fromBotWhatTime:
		tg.pageMarker[message.From.ID] = models.Pages{MarkerWriteTime: true}
	case fromBotScheduleName:
		tg.pageMarker[message.From.ID] = models.Pages{MarkerScheduleName: true}
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
