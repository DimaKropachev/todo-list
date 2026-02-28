package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/DimaKropachev/todo-list/internal/models"
)

type Service interface {
	CreateTask(context.Context, string, string, string) (int64, error)
	GetTasks(context.Context) ([]models.Task, error)
	GetTaskByID(context.Context, int64) (models.Task, error)
	UpdateTask(context.Context, int64, string, string, string) error
	DeleteTask(context.Context, int64) error
}

type Handlers struct {
	s Service
}

func NewHandlers(s Service) *Handlers {
	return &Handlers{
		s: s,
	}
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	t := models.Task{}
	if err = json.Unmarshal(data, &t); err != nil {
		http.Error(w, "invaled json body", http.StatusBadRequest)
		return
	}

	id, err := h.s.CreateTask(ctx, t.Name, t.Desc, t.Status)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.s.GetTasks(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Tasks []models.Task `json:"tasks"`
	}{
		Tasks: tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 0, 64)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task, err := h.s.GetTaskByID(ctx, id)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Task models.Task `json:"task"`
	}{
		Task: task,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()
	id, err := strconv.ParseInt(r.PathValue("id"), 0, 64)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	t := models.Task{}
	if err = json.Unmarshal(data, &t); err != nil {
		http.Error(w, "invaled json body", http.StatusBadRequest)
		return
	}

	err = h.s.UpdateTask(ctx, id, t.Name, t.Desc, t.Status)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 0, 64)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.s.DeleteTask(ctx, id)
	if err != nil {
		// TODO: logging error
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
