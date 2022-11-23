package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
)

const (
	baseGeolocationUrL   = "http://api.openweathermap.org/geo/1.0/direct?q="
	suffixGeolocationUrL = "&limit=5&appid="
	baseWeatherUrl       = "https://api.openweathermap.org/data/2.5/forecast?"
	middleWeatherUrL     = "&appid="
	suffixWeatherUrl     = "&units=metric"
)

type forecast struct {
	weatherForecast models.ForecastResponce
	g               []models.Geolocation
	key             string
}

func newWeather(cfg *models.Config) *forecast {
	return &forecast{key: cfg.WeatherKey}
}

func (tg *telegramBot) setGeolocationRequest(place string) error {
	if err := tg.geolocationUrlValidator(&place); err != nil {
		return err
	}
	url := fmt.Sprint(baseGeolocationUrL, place, suffixGeolocationUrL, tg.forecast.key)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	byteVaule, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll(resp.Body) falied, err:%v", err)
	}
	if err := json.Unmarshal(byteVaule, &tg.forecast.g); err != nil {
		return err
	}
	return nil
}

func (tg *telegramBot) forecastRequest(coordinate string) (string, error) {
	if err := tg.weatherUrlValidator(coordinate); err != nil {
		return "Some problem with forecast. Pleace enter /start again.",
			fmt.Errorf("validator forecast problem. weatherUrlValidator(coordinate) failed err: %v", err)
	}
	url := fmt.Sprint(baseWeatherUrl, coordinate, middleWeatherUrL, tg.forecast.key, suffixWeatherUrl)
	resp, err := http.Get(url)
	if err != nil {
		return "Some problem with forecast. Pleace enter /start again.", fmt.Errorf("http.Get(url) failed, err: %v", err)
	}
	byteVaule, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(byteVaule, &tg.forecast.weatherForecast); err != nil {
		return "Some problem with forecast. Pleace enter /start again.",
			fmt.Errorf("json.Unmarshal(byteVaule, &tg.forecast.w) failed, err: %v", err)
	}
	var res string

	if tg.forecast.weatherForecast.ResponseCode == "200" {
		for count, w := range tg.forecast.weatherForecast.List {
			if err != nil {
				return "", err
			}
			temp := (math.Round(float64(w.Temperature.TemperatureMin+w.Temperature.TemperatureMax) / 2))
			res += fmt.Sprintf("In Time: %v Temperature is %v Humidity is %v\n\n",
				fmt.Sprint(time.Unix(int64(w.DateTime), 0).Format("Mon 15:04:05")), temp,
				w.Temperature.Humidity)
			if count == 10 {
				break
			}
		}
	} else {
		return "Some problem with forecast response. Please enter /start again, or try again later",
			errors.New("Error forecast response")
	}
	return res, nil
}

func (tg *telegramBot) weatherUrlValidator(coordinate string) error {
	if baseWeatherUrl == "" || middleWeatherUrL == "" || suffixWeatherUrl == "" {
		return ErrWeatherUrl
	}
	if coordinate[0] != 'l' || coordinate[1] != 'a' || coordinate[3] != '=' || coordinate[13] != '=' {
		return ErrCoordinate
	}
	if tg.forecast.key == "" {
		return ErrOpenweathermapToken
	}
	return nil
}

func (tg *telegramBot) geolocationUrlValidator(place *string) error {
	if baseGeolocationUrL == "" {
		return ErrbaseGeolocationUrL
	}
	if suffixGeolocationUrL == "" {
		return ErrsuffixGeolocationUrL
	}
	if *place == "" {
		return ErrEmptyPlace
	}
	return nil
}
