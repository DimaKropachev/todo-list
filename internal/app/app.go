package app

import (
	"context"
	"log"

	"github.com/DimaKropachev/todo-list/internal/config"
	"github.com/DimaKropachev/todo-list/internal/repository"
	"github.com/DimaKropachev/todo-list/internal/service"
	"github.com/DimaKropachev/todo-list/internal/transport/http"
	"github.com/DimaKropachev/todo-list/pkg/db"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Start(ctx context.Context) {
	db, err := db.New(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db.DB)
	service := service.New(repo)
	handlers := http.NewHandlers(service)
	router := http.NewRouter(a.cfg.HTTP, handlers)
	
	if err := router.Start(); err != nil  {
		log.Fatal(err)
	}
}
