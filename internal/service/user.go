package service

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
)

func (s *Service) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	hashed, err := auth.HashPassword(u.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	u.Password = hashed

	return s.storer.CreateUser(ctx, u)
}

func (s *Service) GetUser(ctx context.Context, email string) (*model.User, error) {
	return s.storer.GetUser(ctx, email)
}

func (s *Service) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.storer.ListUsers(ctx)
}

func (s *Service) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	if u.Password != "" && !auth.IsHashedPassword(u.Password) {
		hashed, err := auth.HashPassword(u.Password)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %w", err)
		}

		u.Password = hashed
	}

	return s.storer.UpdateUser(ctx, u)
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	return s.storer.DeleteUser(ctx, id)
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.storer.GetUser(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("error getting user: %w", err)
	}

	err = auth.ComparePassword(password, user.Password)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.tokenGen.GenerateAccessToken(
		user.ID,
		user.Email,
		user.IsAdmin,
	)
	if err != nil {
		return nil, "", fmt.Errorf("error generating access token: %w", err)
	}

	return user, accessToken, nil
}
