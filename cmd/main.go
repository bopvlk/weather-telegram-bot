package main

import (
	"fmt"
	"log"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/bot"
)

func run() error {
	var c bot.Container

	// making new configure param
	c, err := bot.NewContainer()
	if err != nil {
		return fmt.Errorf("can not create a new container err: %v", err)
	}

	// making new struct telegrambot with configuratons
	tg, err := bot.New(c)
	if err != nil {
		return fmt.Errorf("New() error. Start application failed %v ", err)
	}

	// running the bot
	tg.Running()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("bot crashed! run().Error : %v", err)
	}
}
