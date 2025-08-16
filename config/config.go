package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
	"github.com/go-playground/validator/v10"
)

type Configuration struct {
	Port             int    `env:"PORT"               envdefault:"8080"   validate:"required,min=1,max=65535"`
	LoggerLevel      string `env:"LOGGER_LEVEL"       envdefault:"INFO"   validate:"required,oneof=DEBUG INFO WARN ERROR"`
	DatabaseURL      string `env:"DATABASE_URL"       validate:"required"`
	ShuttingDowntime int    `env:"SHUTTING_DOWN_TIME" envdefault:"5"      validate:"required,min=5,max=600"`
}

func Load() (*Configuration, error) {
	cfg := &Configuration{}

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
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
