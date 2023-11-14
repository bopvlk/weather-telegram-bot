package main

import (
	"fmt"
	"log"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/bot"
)

func run() error {

	// making new configure param
	container, err := bot.NewContainer()
	if err != nil {
		return fmt.Errorf("can not create a new container err: %v", err)
	}

	// making new struct telegrambot with configuratons
	tg, err := bot.New(container)
	if err != nil {
		return fmt.Errorf("New() error. Start application failed %v ", err)
	}

	// running the bot
	tg.Running()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("bot crashed! run(). error : %v", err)
	}
}
