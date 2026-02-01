package http

import (
	"fmt"
	"net/http"
)

type Config struct {
	Host string `env:"HTTP_HOST"`
	Port int    `env:"HTTP_PORT"`
}

type Router struct {
	cfg Config
	h   *Handlers
}

func NewRouter(cfg Config, h *Handlers) *Router {
	return &Router{
		cfg: cfg,
		h:   h,
	}
}

func (r *Router) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /task", r.h.CreateTask)
	mux.HandleFunc("GET /tasks", r.h.GetTasks)
	mux.HandleFunc("GET /task/{id}", r.h.GetTaskByID)
	mux.HandleFunc("PUT /task/{id}", r.h.UpdateTask)
	mux.HandleFunc("DELETE /task/{id}", r.h.DeleteTask)

	addr := fmt.Sprintf("%s:%d", r.cfg.Host, r.cfg.Port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("failed starting server: %w", err)
	}
	return nil
}
