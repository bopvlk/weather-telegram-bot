package bot

import (
	"errors"
)

var (
	ErrbaseGeolocationUrL   = errors.New("geolocation url is empty")
	ErrsuffixGeolocationUrL = errors.New("suffix part in geolocation url is empty")
	ErrEmptyPlace           = errors.New("place string is empty. ")
	ErrWeatherUrl           = errors.New("some part url in forecast request is empty")
	ErrCoordinate           = errors.New("coordinate string in forecast request is wrong")
	ErrOpenweathermapToken  = errors.New("token from  openweathermap is incorrect")
)
