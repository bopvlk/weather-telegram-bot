package models

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN,required"`
	WeatherKey    string `env:"WEATHER_TOKEN,required"`
	ServiceName   string `env:"SERVICE_NAME" envDefault:"BOT"`
	DBUser        string `env:"MONGO_DB_USER,required"`
	DBPassword    string `env:"MONGO_DB_PASSWORD,required"`
	
}
