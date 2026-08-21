package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (ps *PostgresStorer) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := ps.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back db: %w", rbErr)
		}
		return fmt.Errorf("error in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}