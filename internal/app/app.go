package app

import (
	"context"

	"github.com/DimaKropachev/todo-list/internal/config"
	"github.com/DimaKropachev/todo-list/internal/repository"
	"github.com/DimaKropachev/todo-list/internal/service"
	"github.com/DimaKropachev/todo-list/internal/transport/http"
	"github.com/DimaKropachev/todo-list/pkg/db"
	"github.com/DimaKropachev/todo-list/pkg/logger"
	"go.uber.org/zap"
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
	l := logger.GetLoggerFromCtx(ctx)

	db, err := db.New(a.cfg.DB)
	if err != nil {
		l.Fatal("error to connect to database", zap.Error(err))
	}

	repo := repository.New(db.DB)
	service := service.New(repo)
	handlers := http.NewHandlers(service)
	router := http.NewRouter(a.cfg.HTTP, handlers)

	if err := router.Start(); err != nil {
		l.Fatal("error starting server", zap.Error(err))
	}
}
