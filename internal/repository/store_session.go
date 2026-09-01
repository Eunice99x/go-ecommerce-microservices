package repository

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
)

func (ps *PostgresStorer) CreateSession(ctx context.Context, s *model.Session) (*model.Session, error) {
	query := `INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`

	err := ps.db.GetContext(
		ctx,
		&s.CreatedAt,
		query,
		s.ID,
		s.UserEmail,
		s.RefreshToken,
		s.IsRevoked,
		s.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return s, nil
}

func (ps *PostgresStorer) GetSession(ctx context.Context, id string) (*model.Session, error) {
	var s model.Session

	err := ps.db.GetContext(ctx, &s, "SELECT * FROM sessions WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &s, nil
}

func (ps *PostgresStorer) RevokeSession(ctx context.Context, id string) error {
	_, err := ps.db.ExecContext(ctx, "UPDATE sessions SET is_revoked=true WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (ps *PostgresStorer) DeleteSession(ctx context.Context, id string) error {
	_, err := ps.db.ExecContext(ctx, "DELETE FROM sessions WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}
