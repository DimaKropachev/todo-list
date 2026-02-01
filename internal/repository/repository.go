package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DimaKropachev/todo-list/internal/models"
	"github.com/Masterminds/squirrel"
)

type Repository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Repository) CreateTask(ctx context.Context, name, desc, status string) (int64, error) {
	q, args, err := r.builder.Insert("tasks").Columns("name", "description", "status").Values(name, desc, status).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to create sql query: %w", err)
	}

	var id int64
	if err = r.db.QueryRowContext(ctx, q, args...).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert tast in table: %w", err)
	}

	return id, nil
}

func (r *Repository) GetTasks(ctx context.Context) ([]models.Task, error) {
	q, args, err := r.builder.Select("id", "name", "description", "status").From("tasks").ToSql()
	if err != nil {
		return nil, fmt.Errorf("falied to create sql query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		t := models.Task{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Desc, &t.Status); err != nil {
			return nil, fmt.Errorf("failed to scan result: %w", err)
		}
		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	return tasks, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, id int64) (models.Task, error) {
	q, args, err := r.builder.Select("id", "name", "description", "status").From("tasks").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create sql query: %w", err)
	}

	t := models.Task{}
	row := r.db.QueryRowContext(ctx, q, args...)
	if err = row.Scan(&t.ID, &t.Name, &t.Desc, &t.Status); err != nil {
		return models.Task{}, fmt.Errorf("failed to scan result: %w", err)
	}

	if row.Err() != nil {
		return models.Task{}, fmt.Errorf("failed to get row: %w", err)
	}

	return t, nil
}

func (r *Repository) UpdateTask(ctx context.Context, id int64, name, desc, status string) error {
	qBuilder := r.builder.Update("tasks")
	if name != "" {
		qBuilder = qBuilder.Set("name", name)
	}
	if desc != "" {
		qBuilder = qBuilder.Set("description", desc)
	}
	if status != "" {
		qBuilder = qBuilder.Set("status", status)
	}
	q, args, err := qBuilder.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to create sql query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if n == 0 {
		return fmt.Errorf("no rows deleted (record with id=%d no found)", id)
	}

	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, id int64) error {
	q, args, err := r.builder.Delete("tasks").Where(squirrel.Eq{"id":id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to create sql query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if n == 0 {
		return fmt.Errorf("no rows deleted (record with id=%d not found)", id)
	}

	return nil
}