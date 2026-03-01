package main

import (
	"context"
	"fmt"
	"os"

	"github.com/DimaKropachev/todo-list/internal/app"
	"github.com/DimaKropachev/todo-list/internal/config"
	"github.com/DimaKropachev/todo-list/pkg/logger"
)

func main() {
	path := os.Getenv("ENV_PATH")

	cfg, err := config.ParseConfig(path)
	if err != nil {
		panic(fmt.Errorf("failed parse config: %w", err))
	}

	ctx := context.Background()
	ctx, err = logger.New(ctx, cfg.Env)
	if err != nil {
		panic(fmt.Errorf("failed create logger: %w", err))
	}

	a := app.New(cfg)
	a.Start(ctx)
}
