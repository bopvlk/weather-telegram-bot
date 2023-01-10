package bot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"git.foxminded.com.ua/2.4-weather-forecast-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	updateOffset  = 0
	updateTimeout = 60
)

type telegramBot struct {
	container Container
	bot       *tgbotapi.BotAPI
	updates   tgbotapi.UpdatesChannel
	forecast  *forecast

	store      storage.MongoStorage
	pageMarker map[int64]models.Pages
}

func (tg *telegramBot) Running() {
	// works in a continuous cycle and wait for some events
	for update := range tg.updates {
		tg.eventUpdates(update)
	}
}

func (tg *telegramBot) eventUpdates(update tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	l := tg.container.NewLogger()

	switch {
	case update.CallbackQuery != nil:
		if err := tg.onCallbackQuery(ctx, update.CallbackQuery); err != nil {
			l.Errorf("some error with callback, err: %v", err)
		}
	case update.Message != nil:
		if err := tg.onCommandCreate(ctx, update.Message); err != nil {
			l.Errorf("some error with command, err %v", err)
		}
	default:
		l.Infof("unknown event: %+v\n", update)
	}
}

func New(container Container) (*telegramBot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	l := container.NewLogger()
	cfg := container.NewConfig()

	bot, err := newBot(cfg)
	if err != nil {
		return nil, err
	}

	store, err := storage.NewStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	t := &telegramBot{
		container:  container,
		bot:        bot,
		forecast:   newWeather(cfg),
		store:      store,
		pageMarker: make(map[int64]models.Pages),
	}

	botUpdate := tgbotapi.NewUpdate(updateOffset)
	botUpdate.Timeout = updateTimeout
	t.updates = t.bot.GetUpdatesChan(botUpdate)

	l.Infof("Authorized on account %s", bot.Self.UserName)
	return t, nil
}

func newBot(cfg *models.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(cfg.TelegramToken, tgbotapi.APIEndpoint, &http.Client{Timeout: updateTimeout * time.Second})
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPIWithClient() failed. Error:'%v'\n ", err)
	}
	bot.Debug = true
	return bot, nil
}
