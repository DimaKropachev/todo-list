package config

import (
	"fmt"

	"github.com/DimaKropachev/todo-list/internal/transport/http"
	"github.com/DimaKropachev/todo-list/pkg/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `env:"ENVIRONMENT"`
	DB   db.Config
	HTTP http.Config
}

func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}

	if path != "" {
		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("failed to read env: %w", err)
		}
	}

	return cfg, nil
}
