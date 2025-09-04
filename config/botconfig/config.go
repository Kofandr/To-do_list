package botconfig

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Configuration struct {
	BotToken    string `env:"BOT_TOKEN" validate:"required"`
	BotPort     string `env:"BOT_PORT" envdefault:"8081" validate:"required"`
	APIURL      string `env:"API_URL" envdefault:"http://localhost:8080" validate:"required"`
	LoggerLevel string `env:"LOGGER_LEVEL" envdefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
}

func Load() (*Configuration, error) {
	cfg := &Configuration{}

	godotenv.Load()

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("validator error: %w", err)
	}

	return cfg, nil
}

func MustLoad() *Configuration {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	return cfg
}
