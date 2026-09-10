package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/service"
	"github.com/go-chi/chi/v5"
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
					Password: "hashed-password",
				}

				body := []byte(`{
					"name": "Younes",
					"email": "younes@example.com",
					"password": "plain-password"
				}`)

				req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					user: u,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateUser(rec, req)

				require.Equal(t, http.StatusCreated, rec.Code)

				var res dto.UserRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, int64(1), res.ID)
				require.Equal(t, "Younes", res.Name)
				require.Equal(t, "younes@example.com", res.Email)
				require.False(t, res.IsAdmin)
			},
		},
		{
			name: "failed decoding request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/users",
					bytes.NewBufferString(`{"name":`),
				)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.CreateUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed creating user",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/users",
					bytes.NewReader([]byte(`{"name":"Younes","email":"younes@example.com"}`)),
				)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error creating user"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateUser(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
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

				req := httptest.NewRequest(
					http.MethodGet,
					"/users/user?email=younes@example.com",
					nil,
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					user: u,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.GetUser(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res dto.UserRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, int64(1), res.ID)
				require.Equal(t, "younes@example.com", res.Email)
			},
		},
		{
			name: "missing email",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/users/user", nil)
				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.GetUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed getting user",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodGet,
					"/users/user?email=younes@example.com",
					nil,
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error getting user"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.GetUser(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
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
					{
						ID:    1,
						Name:  "Younes",
						Email: "younes@example.com",
					},
					{
						ID:      2,
						Name:    "Admin",
						Email:   "admin@example.com",
						IsAdmin: true,
					},
				}

				req := httptest.NewRequest(http.MethodGet, "/users", nil)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					users: users,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.ListUsers(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res []dto.UserRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Len(t, res, 2)
				require.Equal(t, "Younes", res[0].Name)
				require.True(t, res[1].IsAdmin)
			},
		},
		{
			name: "failed listing users",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/users", nil)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error listing users"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.ListUsers(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
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
			name: "success",
			test: func(t *testing.T) {
				u := &model.User{
					ID:       1,
					Name:     "Younes",
					Email:    "younes@example.com",
					Password: "hashed-password",
				}

				body := []byte(`{
					"name": "Updated Younes",
					"email": "updated@example.com"
				}`)

				req := httptest.NewRequest(
					http.MethodPatch,
					"/users/user?email=younes@example.com",
					bytes.NewReader(body),
				)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					user: u,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.UpdateUser(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res dto.UserRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, "Updated Younes", res.Name)
				require.Equal(t, "updated@example.com", res.Email)
				require.NotNil(t, res.UpdatedAt)
			},
		},
		{
			name: "missing email",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch,
					"/users/user",
					bytes.NewReader([]byte(`{"name":"Updated Younes"}`)),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.UpdateUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch,
					"/users/user?email=younes@example.com",
					bytes.NewReader([]byte(`{"name":`)),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.UpdateUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed getting user",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch,
					"/users/user?email=younes@example.com",
					bytes.NewReader([]byte(`{"name":"Updated Younes"}`)),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error getting user"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.UpdateUser(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "failed updating user",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch,
					"/users/user?email=younes@example.com",
					bytes.NewReader([]byte(`{"name":"Updated Younes"}`)),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					user:      &model.User{ID: 1, Email: "younes@example.com"},
					updateErr: fmt.Errorf("error updating user"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.UpdateUser(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
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
				req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.DeleteUser(rec, req)

				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name: "invalid user id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.DeleteUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed deleting user",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error deleting user"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.DeleteUser(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
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
				now := time.Now()

				result := &service.LoginResult{
					User: &model.User{
						ID:    1,
						Name:  "Younes",
						Email: "younes@example.com",
					},
					SessionID:             "8d2a6f1e-0e2f-4f3a-9a1e-2b3c4d5e6f70",
					AccessToken:           "access-token",
					RefreshToken:          "refresh-token",
					AccessTokenExpiresAt:  now.Add(15 * time.Minute),
					RefreshTokenExpiresAt: now.Add(7 * 24 * time.Hour),
				}

				body := []byte(`{
					"email": "younes@example.com",
					"password": "plain-password"
				}`)

				req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					loginResult: result,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.LoginUser(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res dto.LoginUserRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, result.SessionID, res.SessionID)
				require.Equal(t, "access-token", res.AccessToken)
				require.Equal(t, "refresh-token", res.RefreshToken)
				require.Equal(t, int64(1), res.User.ID)
				require.Equal(t, "younes@example.com", res.User.Email)
			},
		},
		{
			name: "failed decoding request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/login",
					bytes.NewBufferString(`{"email":`),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.LoginUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid credentials",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/login",
					bytes.NewReader([]byte(`{"email":"younes@example.com","password":"wrong"}`)),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("invalid credentials"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.LoginUser(rec, req)

				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestLogoutUser(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/logout/session-id", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "session-id")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.LogoutUser(rec, req)

				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name: "missing session id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/logout/", nil)

				rctx := chi.NewRouteContext()

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.LogoutUser(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed deleting session",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/logout/session-id", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "session-id")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error deleting session"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.LogoutUser(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
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
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				expiresAt := time.Now().Add(15 * time.Minute)

				req := httptest.NewRequest(
					http.MethodPost,
					"/refresh",
					bytes.NewReader([]byte(`{"refresh_token":"refresh-token"}`)),
				)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					accessToken: "new-access-token",
					expiresAt:   expiresAt,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.RenewAccessToken(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res dto.RenewAccessTokenRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, "new-access-token", res.AccessToken)
				require.WithinDuration(t, expiresAt, res.AccessTokenExpiresAt, time.Second)
			},
		},
		{
			name: "failed decoding request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/refresh",
					bytes.NewBufferString(`{"refresh_token":`),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.RenewAccessToken(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid refresh token",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/refresh",
					bytes.NewReader([]byte(`{"refresh_token":"expired-token"}`)),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("invalid refresh token"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.RenewAccessToken(rec, req)

				require.Equal(t, http.StatusUnauthorized, rec.Code)
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
				req := httptest.NewRequest(
					http.MethodPatch,
					"/sessions/session-id/revoke",
					nil,
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "session-id")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.RevokeSession(rec, req)

				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name: "missing session id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPatch, "/sessions//revoke", nil)

				rctx := chi.NewRouteContext()

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.RevokeSession(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed revoking session",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch,
					"/sessions/session-id/revoke",
					nil,
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "session-id")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error revoking session"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.RevokeSession(rec, req)

				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestPatchUserReq(t *testing.T) {
	base := func() *model.User {
		return &model.User{
			ID:       1,
			Name:     "Younes",
			Email:    "younes@example.com",
			Password: "hashed-password",
		}
	}

	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "patches every field",
			test: func(t *testing.T) {
				u := base()

				patchUserReq(u, dto.UserReq{
					Name:     "Updated Younes",
					Email:    "updated@example.com",
					Password: "new-password",
				})

				require.Equal(t, "Updated Younes", u.Name)
				require.Equal(t, "updated@example.com", u.Email)
				require.Equal(t, "new-password", u.Password)
				require.NotNil(t, u.UpdatedAt)
			},
		},
		{
			name: "leaves omitted fields untouched",
			test: func(t *testing.T) {
				u := base()

				patchUserReq(u, dto.UserReq{Name: "Updated Younes"})

				require.Equal(t, "Updated Younes", u.Name)
				require.Equal(t, "younes@example.com", u.Email)
				require.Equal(t, "hashed-password", u.Password)
			},
		},
		{
			name: "an empty request only stamps updated_at",
			test: func(t *testing.T) {
				u := base()
				want := base()

				patchUserReq(u, dto.UserReq{})

				require.NotNil(t, u.UpdatedAt)

				u.UpdatedAt = nil
				require.Equal(t, want, u)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
