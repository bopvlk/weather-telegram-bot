package bot

import (
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// load configuration information of bot
func SetUpConfig() (*models.Config, error) {
	config := &models.Config{}

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("godotenv.load failed. error:'%v'", err)
	}

	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("env.Parse() in bot.config failed. Error:'%v'", err)
	}

	return config, nil
}
