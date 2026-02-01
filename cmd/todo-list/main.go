package main

import (
	"context"
	"log"
	"os"

	"github.com/DimaKropachev/todo-list/internal/app"
	"github.com/DimaKropachev/todo-list/internal/config"
)

func main() {
	path := os.Getenv("ENV_PATH")

	cfg, err := config.ParseConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	a := app.New(cfg)
	a.Start(ctx)
}
