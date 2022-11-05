package models

type Config struct {
	TelegramToken string `env:"TOKEN,required"`
	WeatherKey    string `env:"WEATHER,required"`
	ServiceName   string `env:"SERVICE" envDefault:""`
}
