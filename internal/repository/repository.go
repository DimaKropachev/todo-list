package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type Repository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

