package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	w   forecastResponce
	g   []geolocation
	key string
}

type geolocation struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
	Country   string  `json:"country"`
	Region    string  `json:"state"`
}

type forecastResponce struct {
	ResponseCode string `json:"cod"`
	List         []struct {
		DateTime    int `json:"dt"`
		Temperature struct {
			TemperatureMin float32 `json:"temp_min"`
			TemperatureMax float32 `json:"temp_max"`
			Humidity       int     `json:"humidity"`
		} `json:"main"`
		SkyWeather []struct {
			InSky          string `json:"main"`
			DescriptionSky string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
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
	byteVaule, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(byteVaule, &tg.forecast.g); err != nil {
		return err
	}
	return nil
}

func (tg *telegramBot) forecastRequest(coordinate string) (string, error) {
	if err := tg.weatherUrlValidator(coordinate); err != nil {
		return "Some problem with forecast. Pleace enter /start again.", fmt.Errorf("validator forecast problem. weatherUrlValidator(coordinate) failed err: %v", err)
	}
	url := fmt.Sprint(baseWeatherUrl, coordinate, middleWeatherUrL, tg.forecast.key, suffixWeatherUrl)
	resp, err := http.Get(url)
	if err != nil {
		return "Some problem with forecast. Pleace enter /start again.", fmt.Errorf("http.Get(url) failed, err: %v", err)
	}
	byteVaule, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(byteVaule, &tg.forecast.w); err != nil {
		return "Some problem with forecast. Pleace enter /start again.",
			fmt.Errorf("json.Unmarshal(byteVaule, &tg.forecast.w) failed, err: %v", err)
	}
	var res string

	if tg.forecast.w.ResponseCode == "200" {
		for count, w := range tg.forecast.w.List {
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
