package models

type Geolocation struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
	Country   string  `json:"country"`
	Region    string  `json:"state"`
}

type Temperature struct {
	TemperatureMin float32 `json:"temp_min"`
	TemperatureMax float32 `json:"temp_max"`
	Humidity       int     `json:"humidity"`
}

type SkyWeather struct {
	InSky          string `json:"main"`
	DescriptionSky string `json:"description"`
}

type ForecastResponce struct {
	ResponseCode string `json:"cod"`
	List         []struct {
		DateTime    int          `json:"dt"`
		Temperature Temperature  `json:"main"`
		SkyWeathers []SkyWeather `json:"weather"`
	} `json:"list"`
}
