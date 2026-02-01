package http

import "net/http"

type Config struct {
	Host string
	Port int
}

type Router struct {
	h *Handlers
}

func NewRouter(h *Handlers) *Router {
	return &Router{
		h: h,
	}
}

func (r *Router) Start() error {
	mux := http.NewServeMux()

	
}