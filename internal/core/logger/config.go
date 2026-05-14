package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string  `envconfig:"LEVEL" required:"true"`
	Folder string  `envconfig:"FOLDER" required:"true"`
}

//читает env, заполняет struct 
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("proccess envconfig: %w", err)
	}

	return config, nil
}

//либо конфиг есть, либо приложение вообще не должно запускаться
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}
	return config
}