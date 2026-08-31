package repository

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
)

func (ps *PostgresStorer) CreateUser(ctx context.Context, u *model.User) (*model.User, error){
	query := `INSERT INTO users (name, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := ps.db.GetContext(ctx, u, query, u.Name, u.Email, u.Password, u.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return u, nil
}

func (ps *PostgresStorer) GetUser(ctx context.Context, email string) (*model.User, error) {
	var u model.User

	err := ps.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &u, nil
}

func (ps *PostgresStorer) ListUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	err := ps.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (ps *PostgresStorer) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET
			name=$1,
			email=$2,
			password=$3,
			is_admin=$4,
			updated_at=$5
		WHERE id=$6
	`

	_, err := ps.db.ExecContext(
		ctx,
		query,
		u.Name,
		u.Email,
		u.Password,
		u.IsAdmin,
		u.UpdatedAt,
		u.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

func (ps *PostgresStorer) DeleteUser(ctx context.Context, id int64) error {
	_, err := ps.db.ExecContext(
		ctx,
		"DELETE FROM users WHERE id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}