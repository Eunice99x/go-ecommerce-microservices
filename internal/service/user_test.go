package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				u := &model.User{
					ID:       1,
					Name:     "Younes",
					Email:    "younes@example.com",
					Password: "plain-password",
				}

				s := &Service{
					storer: &fakeStorer{user: u},
				}

				got, err := s.CreateUser(t.Context(), u)

				require.NoError(t, err)
				require.Equal(t, int64(1), got.ID)
				require.NotEqual(t, "plain-password", got.Password)
				require.True(t, auth.IsHashedPassword(got.Password))
				require.NoError(t, auth.ComparePassword("plain-password", got.Password))
			},
		},
		{
			name: "failed creating user",
			test: func(t *testing.T) {
				u := &model.User{
					Name:     "Younes",
					Email:    "younes@example.com",
					Password: "plain-password",
				}

				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error creating user")},
				}

				_, err := s.CreateUser(t.Context(), u)

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

func TestGetUser(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				u := &model.User{
					ID:    1,
					Name:  "Younes",
					Email: "younes@example.com",
				}

				s := &Service{
					storer: &fakeStorer{user: u},
				}

				got, err := s.GetUser(t.Context(), "younes@example.com")

				require.NoError(t, err)
				require.Equal(t, int64(1), got.ID)
				require.Equal(t, "younes@example.com", got.Email)
			},
		},
		{
			name: "failed getting user",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error getting user")},
				}

				_, err := s.GetUser(t.Context(), "younes@example.com")

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

func TestListUsers(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				users := []*model.User{
					{ID: 1, Name: "Younes", Email: "younes@example.com"},
					{ID: 2, Name: "Admin", Email: "admin@example.com", IsAdmin: true},
				}

				s := &Service{
					storer: &fakeStorer{users: users},
				}

				got, err := s.ListUsers(t.Context())

				require.NoError(t, err)
				require.Len(t, got, 2)
				require.Equal(t, "Younes", got[0].Name)
				require.True(t, got[1].IsAdmin)
			},
		},
		{
			name: "failed listing users",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error listing users")},
				}

				_, err := s.ListUsers(t.Context())

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

func TestUpdateUser(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "hashes a new plain password",
			test: func(t *testing.T) {
				u := &model.User{
					ID:       1,
					Name:     "Younes",
					Email:    "younes@example.com",
					Password: "new-plain-password",
				}

				s := &Service{
					storer: &fakeStorer{user: u},
				}

				got, err := s.UpdateUser(t.Context(), u)

				require.NoError(t, err)
				require.True(t, auth.IsHashedPassword(got.Password))
				require.NoError(t, auth.ComparePassword("new-plain-password", got.Password))
			},
		},
		{
			name: "keeps an already hashed password",
			test: func(t *testing.T) {
				hashed, err := auth.HashPassword("plain-password")
				require.NoError(t, err)

				u := &model.User{
					ID:       1,
					Name:     "Younes",
					Email:    "younes@example.com",
					Password: hashed,
				}

				s := &Service{
					storer: &fakeStorer{user: u},
				}

				got, err := s.UpdateUser(t.Context(), u)

				require.NoError(t, err)
				require.Equal(t, hashed, got.Password)
			},
		},
		{
			name: "keeps an empty password",
			test: func(t *testing.T) {
				u := &model.User{
					ID:    1,
					Name:  "Younes",
					Email: "younes@example.com",
				}

				s := &Service{
					storer: &fakeStorer{user: u},
				}

				got, err := s.UpdateUser(t.Context(), u)

				require.NoError(t, err)
				require.Empty(t, got.Password)
			},
		},
		{
			name: "failed updating user",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error updating user")},
				}

				_, err := s.UpdateUser(t.Context(), &model.User{ID: 1})

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

func TestDeleteUser(t *testing.T) {
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

				err := s.DeleteUser(t.Context(), 1)

				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting user",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error deleting user")},
				}

				err := s.DeleteUser(t.Context(), 1)

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

func TestLoginUser(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				hashed, err := auth.HashPassword("plain-password")
				require.NoError(t, err)

				u := &model.User{
					ID:       1,
					Name:     "Younes",
					Email:    "younes@example.com",
					Password: hashed,
				}

				session := &model.Session{
					ID:        "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70",
					UserEmail: u.Email,
				}

				tokenGen := auth.DefaultJWTConfig("secret")

				s := &Service{
					storer:   &fakeStorer{user: u, session: session},
					tokenGen: tokenGen,
				}

				got, err := s.LoginUser(t.Context(), u.Email, "plain-password")

				require.NoError(t, err)
				require.Equal(t, u, got.User)
				require.Equal(t, session.ID, got.SessionID)
				require.NotEmpty(t, got.AccessToken)
				require.NotEmpty(t, got.RefreshToken)
				require.True(t, got.RefreshTokenExpiresAt.After(got.AccessTokenExpiresAt))

				accessClaims, err := tokenGen.ValidateToken(got.AccessToken, "access")
				require.NoError(t, err)
				require.Equal(t, u.ID, accessClaims.ID)
				require.Equal(t, u.Email, accessClaims.Email)

				refreshClaims, err := tokenGen.ValidateToken(got.RefreshToken, "refresh")
				require.NoError(t, err)
				require.Equal(t, u.Email, refreshClaims.Email)
			},
		},
		{
			name: "failed getting user",
			test: func(t *testing.T) {
				s := &Service{
					storer:   &fakeStorer{err: fmt.Errorf("error getting user")},
					tokenGen: auth.DefaultJWTConfig("secret"),
				}

				_, err := s.LoginUser(t.Context(), "younes@example.com", "plain-password")

				require.Error(t, err)
			},
		},
		{
			name: "wrong password",
			test: func(t *testing.T) {
				hashed, err := auth.HashPassword("plain-password")
				require.NoError(t, err)

				u := &model.User{
					ID:       1,
					Email:    "younes@example.com",
					Password: hashed,
				}

				s := &Service{
					storer:   &fakeStorer{user: u},
					tokenGen: auth.DefaultJWTConfig("secret"),
				}

				_, err = s.LoginUser(t.Context(), u.Email, "wrong-password")

				require.Error(t, err)
			},
		},
		{
			name: "failed creating session",
			test: func(t *testing.T) {
				hashed, err := auth.HashPassword("plain-password")
				require.NoError(t, err)

				u := &model.User{
					ID:       1,
					Email:    "younes@example.com",
					Password: hashed,
				}

				s := &Service{
					storer: &fakeStorer{
						user:             u,
						createSessionErr: fmt.Errorf("error creating session"),
					},
					tokenGen: auth.DefaultJWTConfig("secret"),
				}

				_, err = s.LoginUser(t.Context(), u.Email, "plain-password")

				require.Error(t, err)
			},
		},
		{
			name: "failed validating refresh token",
			test: func(t *testing.T) {
				hashed, err := auth.HashPassword("plain-password")
				require.NoError(t, err)

				u := &model.User{
					ID:       1,
					Email:    "younes@example.com",
					Password: hashed,
				}

				// a refresh token minted already expired fails the refresh
				// claims lookup
				tokenGen := &auth.JWTConfig{
					SecretKey:          "secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: -time.Hour,
					Issuer:             "go-ecommerce",
				}

				s := &Service{
					storer:   &fakeStorer{user: u},
					tokenGen: tokenGen,
				}

				_, err = s.LoginUser(t.Context(), u.Email, "plain-password")

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
