package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
)

func (s *Service) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	createdSession, err := s.storer.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	return createdSession, nil
}

func (s *Service) GetSession(ctx context.Context, id string) (*model.Session, error) {
	session, err := s.storer.GetSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return session, nil
}

func (s *Service) RevokeSession(ctx context.Context, id string) error {
	if err := s.storer.RevokeSession(ctx, id); err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (s *Service) DeleteSession(ctx context.Context, id string) error {
	if err := s.storer.DeleteSession(ctx, id); err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}

func (s *Service) RenewAccessToken(ctx context.Context, refreshToken string) (string, time.Time, error) {
	claims, err := s.tokenGen.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return "", time.Time{}, fmt.Errorf("%w: invalid refresh token: %w", model.ErrInvalidCredentials, err)
	}

	session, err := s.storer.GetSession(
		ctx,
		claims.RegisteredClaims.ID,
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error getting session: %w", err)
	}

	if session.IsRevoked {
		return "", time.Time{}, fmt.Errorf("%w: session revoked", model.ErrInvalidCredentials)
	}

	if time.Now().After(session.ExpiresAt) {
		return "", time.Time{}, fmt.Errorf("%w: session expired", model.ErrInvalidCredentials)
	}

	if session.UserEmail != claims.Email {
		return "", time.Time{}, fmt.Errorf("%w: invalid session", model.ErrInvalidCredentials)
	}

	if session.RefreshToken != refreshToken {
		return "", time.Time{}, fmt.Errorf("%w: invalid refresh token", model.ErrInvalidCredentials)
	}

	accessToken, err := s.tokenGen.GenerateAccessToken(
		claims.ID,
		claims.Email,
		claims.IsAdmin,
		session.ID,
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf(
			"error generating access token: %w",
			err,
		)
	}

	expiresAt := time.Now().Add(s.tokenGen.AccessTokenExpiry)

	return accessToken, expiresAt, nil
}
