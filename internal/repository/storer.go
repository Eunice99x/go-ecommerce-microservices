package repository

import (
	"github.com/jmoiron/sqlx"
)

type PostgresStorer struct {
	db *sqlx.DB
}

func NewPostgresStorer(db *sqlx.DB) *PostgresStorer {
	return &PostgresStorer{
		db: db,
	}
}
