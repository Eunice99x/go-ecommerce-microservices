package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				session := &model.Session{
					ID:           "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70",
					UserEmail:    "younes@example.com",
					RefreshToken: "refresh-token",
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				}

				s := &Service{
					storer: &fakeStorer{session: session},
				}

				got, err := s.CreateSession(t.Context(), session)

				require.NoError(t, err)
				require.Equal(t, session.ID, got.ID)
				require.Equal(t, "younes@example.com", got.UserEmail)
			},
		},
		{
			name: "failed creating session",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error creating session")},
				}

				_, err := s.CreateSession(t.Context(), &model.Session{})

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestGetSession(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				session := &model.Session{
					ID:        "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70",
					UserEmail: "younes@example.com",
				}

				s := &Service{
					storer: &fakeStorer{session: session},
				}

				got, err := s.GetSession(t.Context(), session.ID)

				require.NoError(t, err)
				require.Equal(t, session.ID, got.ID)
			},
		},
		{
			name: "failed getting session",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error getting session")},
				}

				_, err := s.GetSession(t.Context(), "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70")

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestRevokeSession(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{},
				}

				err := s.RevokeSession(t.Context(), "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70")

				require.NoError(t, err)
			},
		},
		{
			name: "failed revoking session",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error revoking session")},
				}

				err := s.RevokeSession(t.Context(), "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70")

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestDeleteSession(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{},
				}

				err := s.DeleteSession(t.Context(), "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70")

				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting session",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error deleting session")},
				}

				err := s.DeleteSession(t.Context(), "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70")

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestRenewAccessToken(t *testing.T) {
	// builds a refresh token plus the session the service should find for it
	newRefreshToken := func(t *testing.T, c *auth.JWTConfig, u *model.User) (string, *model.Session) {
		t.Helper()

		token, err := c.GenerateRefreshToken(u.ID, u.Email, u.IsAdmin)
		require.NoError(t, err)

		claims, err := c.ValidateToken(token, "refresh")
		require.NoError(t, err)

		return token, &model.Session{
			ID:           claims.RegisteredClaims.ID,
			UserEmail:    u.Email,
			RefreshToken: token,
			IsRevoked:    false,
			ExpiresAt:    time.Now().Add(c.RefreshTokenExpiry),
		}
	}

	u := &model.User{
		ID:      1,
		Email:   "younes@example.com",
		IsAdmin: false,
	}

	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				refreshToken, session := newRefreshToken(t, tokenGen, u)

				s := &Service{
					storer:   &fakeStorer{session: session},
					tokenGen: tokenGen,
				}

				accessToken, expiresAt, err := s.RenewAccessToken(t.Context(), refreshToken)

				require.NoError(t, err)
				require.NotEmpty(t, accessToken)
				require.True(t, expiresAt.After(time.Now()))

				claims, err := tokenGen.ValidateToken(accessToken, "access")
				require.NoError(t, err)
				require.Equal(t, u.ID, claims.ID)
				require.Equal(t, u.Email, claims.Email)
			},
		},
		{
			name: "invalid refresh token",
			test: func(t *testing.T) {
				s := &Service{
					storer:   &fakeStorer{},
					tokenGen: auth.DefaultJWTConfig("secret"),
				}

				_, _, err := s.RenewAccessToken(t.Context(), "not-a-token")

				require.Error(t, err)
			},
		},
		{
			name: "access token is not accepted as a refresh token",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				accessToken, err := tokenGen.GenerateAccessToken(u.ID, u.Email, u.IsAdmin)
				require.NoError(t, err)

				s := &Service{
					storer:   &fakeStorer{},
					tokenGen: tokenGen,
				}

				_, _, err = s.RenewAccessToken(t.Context(), accessToken)

				require.Error(t, err)
			},
		},
		{
			name: "failed getting session",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				refreshToken, _ := newRefreshToken(t, tokenGen, u)

				s := &Service{
					storer:   &fakeStorer{err: fmt.Errorf("error getting session")},
					tokenGen: tokenGen,
				}

				_, _, err := s.RenewAccessToken(t.Context(), refreshToken)

				require.Error(t, err)
			},
		},
		{
			name: "session revoked",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				refreshToken, session := newRefreshToken(t, tokenGen, u)
				session.IsRevoked = true

				s := &Service{
					storer:   &fakeStorer{session: session},
					tokenGen: tokenGen,
				}

				_, _, err := s.RenewAccessToken(t.Context(), refreshToken)

				require.Error(t, err)
			},
		},
		{
			name: "session expired",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				refreshToken, session := newRefreshToken(t, tokenGen, u)
				session.ExpiresAt = time.Now().Add(-time.Hour)

				s := &Service{
					storer:   &fakeStorer{session: session},
					tokenGen: tokenGen,
				}

				_, _, err := s.RenewAccessToken(t.Context(), refreshToken)

				require.Error(t, err)
			},
		},
		{
			name: "session belongs to another user",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				refreshToken, session := newRefreshToken(t, tokenGen, u)
				session.UserEmail = "someone@example.com"

				s := &Service{
					storer:   &fakeStorer{session: session},
					tokenGen: tokenGen,
				}

				_, _, err := s.RenewAccessToken(t.Context(), refreshToken)

				require.Error(t, err)
			},
		},
		{
			name: "session holds a different refresh token",
			test: func(t *testing.T) {
				tokenGen := auth.DefaultJWTConfig("secret")

				refreshToken, session := newRefreshToken(t, tokenGen, u)
				session.RefreshToken = "another-refresh-token"

				s := &Service{
					storer:   &fakeStorer{session: session},
					tokenGen: tokenGen,
				}

				_, _, err := s.RenewAccessToken(t.Context(), refreshToken)

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
