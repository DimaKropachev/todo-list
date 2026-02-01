package service

import (
	"context"

	"github.com/DimaKropachev/todo-list/internal/models"
)

type Repository interface {
	CreateTask(context.Context, string, string, string) (int64, error)
	GetTasks(context.Context) ([]models.Task, error)
	GetTaskByID(context.Context, int64) (models.Task, error)
	UpdateTask(context.Context, int64, string, string, string) error
	DeleteTask(context.Context, int64) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, name, desc, status string) (int64, error) {
	return s.repo.CreateTask(ctx, name, desc, status)
}

func (s *Service) GetTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetTasks(ctx)
}

func (s *Service) GetTaskByID(ctx context.Context, id int64) (models.Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}

func (s *Service) UpdateTask(ctx context.Context, id int64, name, desc, status string) error {
	return s.repo.UpdateTask(ctx, id, name, desc, status)
}

func (s *Service) DeleteTask(ctx context.Context, id int64) error {
	return s.repo.DeleteTask(ctx, id)
}