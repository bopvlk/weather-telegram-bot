package bot

import (
	"fmt"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
)

type Container interface {
	NewConfig() *models.Config
	NewLogger() *LogrusLogger
}

type container struct {
	config *models.Config
	logger *LogrusLogger
}

func NewContainer() (*container, error) {
	conf, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("NewConfig() error. Start application failed %v ", err)
	}

	// create logger
	l, err := NewLogger(conf)
	if err != nil {
		return nil, fmt.Errorf("Something wrong with logger: %v", err)
	}
	return &container{
		config: conf,
		logger: l,
	}, nil
}

func (c *container) NewConfig() *models.Config {
	return c.config
}

func (c *container) NewLogger() *LogrusLogger {
	return c.logger
}
