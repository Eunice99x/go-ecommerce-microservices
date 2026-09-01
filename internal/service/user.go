package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
)

type LoginResult struct {
	User                  *model.User
	SessionID             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

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

func (s *Service) LoginUser(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.storer.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	err = auth.ComparePassword(password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.tokenGen.GenerateAccessToken(
		user.ID,
		user.Email,
		user.IsAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := s.tokenGen.GenerateRefreshToken(
		user.ID,
		user.Email,
		user.IsAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	refreshClaims, err := s.tokenGen.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return nil, fmt.Errorf("error validating refresh token: %w", err)
	}

	now := time.Now()

	accessTokenExpiresAt := now.Add(s.tokenGen.AccessTokenExpiry)
	refreshTokenExpiresAt := now.Add(s.tokenGen.RefreshTokenExpiry)

	session := &model.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshTokenExpiresAt,
	}

	session, err = s.storer.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return &LoginResult{
		User:                  user,
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}
