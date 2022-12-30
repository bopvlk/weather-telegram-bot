package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jasonlvhit/gocron"
)

func (tg *telegramBot) onCommandCreate(ctx context.Context, message *tgbotapi.Message) error {
	switch {
	case message.Text == "/start":
		user, err := tg.storage.FindUsersPerTelegramId(ctx, message.From.ID)
		if err != nil {
			return err
		}
		if user == nil {
			if err := tg.printMessage(message, fromBotStartIfNoLogged); err != nil {
				return err
			}
		} else {
			if err := tg.printMessage(message, fromBotStartIfLogged); err != nil {
				return err
			}
		}
	case markerWriteTime:
		if timeChecker(message.Text) {
			markerWriteTime = false
			toDBEventTime = message.Text
			if err := tg.printMessage(message, fromBotScheduleName); err != nil {
				return err
			}
		} else {
			if err := tg.printMessage(message, fromBotWrongFormatTime); err != nil {
				return err
			}
		}
	case markerScheduleName:
		markerScheduleName = false
		if _, err := tg.storage.SaveEvent(ctx, toDBEventTime, message.Text); err != nil {
			return err
		}

		scheduleRun := func() error {
			forecast, err := tg.forecastRequest(tg.storage.User.City)
			if err != nil {
				return err
			}
			if err := tg.printMessage(message, forecast); err != nil {
				return err
			}
			return nil
		}

		for _, e := range tg.storage.Events {
			if err := gocron.Every(1).Day().At(e.EventTime).Do(scheduleRun); err != nil {
				return err
			}
		}

		<-gocron.Start()

		if err := tg.printMessage(message, "Now you will get notifications with forecast"); err != nil {
			return err
		}
	case markerFindCity:
		markerFindCity = false
		if err := tg.printMessage(message, fromBotSelectPlace); err != nil {
			return err
		}
	case markerWritePassword:
		markerWritePassword = false
		var err error
		toDBPasswordHash, err = middleware.JwtHashing(message.Text, message.From.ID)
		if err != nil {
			return err
		}
		if err := tg.printMessage(message, fromBotWitchCity); err != nil {
			return err
		}
	default:
		if err := tg.printMessage(message, fromBotUknownMessage); err != nil {
			return err
		}
	}
	return nil
}

func (tg *telegramBot) onCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	clbck := tgbotapi.NewCallback(callback.ID, callback.Data)
	if _, err := tg.bot.Request(clbck); err != nil {
		return fmt.Errorf("send request failed: %v", err)
	}
	switch clbck.Text {
	case keyStartCheckOut:
		if err := tg.printMessage(callback.Message, fromBotWriteCity); err != nil {
			return err
		}
	case keyStartAuthorization:
		if err := tg.printMessage(callback.Message, fromBotPasword); err != nil {
			return err
		}
	case keyJustLoggedYES:
		user, err := tg.storage.FindUsersPerTelegramId(ctx, callback.From.ID)
		if err != nil {
			return err
		}
		forecast, err := tg.forecastRequest(user.City)
		if err != nil {
			return err
		}
		if err := tg.printMessage(callback.Message, forecast); err != nil {
			return err
		}
	case keyJustLoggedSchedule:
		if err := tg.printMessage(callback.Message, fromBotWhatTime); err != nil {
			return err
		}
	}
	if coordinateChecker(clbck.Text) {
		forecast, err := tg.forecastRequest(clbck.Text)
		if err != nil {
			return err
		}
		if markerSaveCityMarker {
			markerSaveCityMarker = false
			if _, err := tg.storage.SaveUser(ctx, callback.From.ID, toDBPasswordHash, clbck.Text); err != nil {
				return err
			}
			toDBPasswordHash = ""
			if err := tg.printMessage(callback.Message, fromBotULogged); err != nil {
				return err
			}
		} else {
			if err := tg.printMessage(callback.Message, forecast); err != nil {
				return err
			}
		}
	}
	return nil
}

func timeChecker(t string) bool {
	hourStr, minuteStr, found := strings.Cut(t, ":")
	if !found {
		return false
	}
	hourInt, err := strconv.Atoi(hourStr)
	if err != nil {
		return false
	}
	if hourInt > 23 || hourInt < 0 {
		return false
	}
	minuteInt, err := strconv.Atoi(minuteStr)
	if err != nil {
		return false
	}
	if minuteInt > 59 || minuteInt < 0 {
		return false
	}
	return true
}

func coordinateChecker(coordinate string) bool {
	if coordinate[0] == 'l' && coordinate[1] == 'a' && coordinate[2] == 't' {
		return true
	}
	return false
}
