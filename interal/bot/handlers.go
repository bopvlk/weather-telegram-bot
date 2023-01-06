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
		user, err := tg.storage.NewUserRepository().FindUser(ctx, message.From.ID)
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
	case tg.pageMarker[message.From.ID].MarkerWriteTime:
		if err := timeChecker(message.Text); err != nil {
			if err := tg.printMessage(message, fromBotWrongFormatTime); err != nil {
				return err
			} else {
				delete(tg.pageMarker, message.From.ID)
				toDBEventTime = message.Text
				if err := tg.printMessage(message, fromBotScheduleName); err != nil {
					return err
				}
			}
		}
	case tg.pageMarker[message.From.ID].MarkerScheduleName:
		delete(tg.pageMarker, message.From.ID)
		user, err := tg.storage.NewUserRepository().FindUser(ctx, message.From.ID)
		if err != nil {
			return err
		}
		if _, err := tg.storage.NewEventRepository().SaveEvent(ctx, user.ID, toDBEventTime, message.Text); err != nil {
			return err
		}

		scheduleRun := func() error {
			forecast, err := tg.forecastRequest(user.City)
			if err != nil {
				return err
			}
			if err := tg.printMessage(message, forecast); err != nil {
				return err
			}
			return nil
		}

		events, err := tg.storage.NewEventRepository().FindEvents(ctx, user.ID)
		if err != nil {
			return err
		}
		for _, e := range events {
			if err := gocron.Every(1).Day().At(e.EventTime).Do(scheduleRun); err != nil {
				return err
			}
		}
		<-gocron.Start()

		if err := tg.printMessage(message, "Now you will get notifications with forecast"); err != nil {
			return err
		}
	case tg.pageMarker[message.From.ID].MarkerFindCity:
		delete(tg.pageMarker, message.From.ID)
		if err := tg.printMessage(message, fromBotSelectPlace); err != nil {
			return err
		}
	case tg.pageMarker[message.From.ID].MarkerWritePassword:
		delete(tg.pageMarker, message.From.ID)
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
		user, err := tg.storage.NewUserRepository().FindUser(ctx, callback.From.ID)
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
		if tg.pageMarker[callback.From.ID].MarkerSaveCityMarker {
			delete(tg.pageMarker, callback.From.ID)
			if _, err := tg.storage.NewUserRepository().SaveUser(ctx, callback.From.ID, toDBPasswordHash, clbck.Text); err != nil {
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

func timeChecker(t string) error {
	hourStr, minuteStr, found := strings.Cut(t, ":")
	if !found {
		return ErrCutProblem
	}
	hourInt, err := strconv.Atoi(hourStr)
	if err != nil {
		return fmt.Errorf("some probem with convert hour from string to int format in timeChecker(). err -> %v", err)
	}
	if hourInt > 23 || hourInt < 0 {
		return ErrWrongHourNum
	}
	minuteInt, err := strconv.Atoi(minuteStr)
	if err != nil {
		return fmt.Errorf("some probem with convert minute from string to int format in timeChecker(). err -> %v", err)
	}
	if minuteInt > 59 || minuteInt < 0 {
		return ErrWrongMinNum
	}
	return nil
}

func coordinateChecker(coordinate string) bool {
	if coordinate[0] == 'l' && coordinate[1] == 'a' && coordinate[2] == 't' {
		return true
	}
	return false
}
