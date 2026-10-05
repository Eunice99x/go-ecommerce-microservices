package repository

import (
	"database/sql"
	"errors"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/lib/pq"
)

const pqUniqueViolation = "23505"

// dbError keeps database/sql and pq errors inside the repository
func dbError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == pqUniqueViolation {
		return model.ErrAlreadyExists
	}

	return err
}
