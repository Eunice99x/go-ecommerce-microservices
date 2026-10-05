package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/go-chi/chi/v5"
)

func toUserModel(u dto.UserReq) *model.User {
	return &model.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func toUserRes(u *model.User) dto.UserRes {
	return dto.UserRes{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func patchUserReq(u *model.User, req dto.UserReq) {
	if req.Name != "" {
		u.Name = req.Name
	}

	if req.Email != "" {
		u.Email = req.Email
	}

	if req.Password != "" {
		u.Password = req.Password
	}

	now := time.Now()

	u.UpdatedAt = &now
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), toUserModel(req))
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}

	res := toUserRes(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// users get themselves, admins may look up anyone by ?email=
	email := claims.Email
	if q := r.URL.Query().Get("email"); q != "" && claims.IsAdmin {
		email = q
	}

	user, err := h.service.GetUser(r.Context(), email)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}

	res := toUserRes(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}

	res := make([]dto.UserRes, 0, len(users))

	for _, user := range users {
		res = append(res, toUserRes(user))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	u, err := h.service.GetUser(r.Context(), claims.Email)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}

	patchUserReq(u, req)
	if u.Email == "" {
		u.Email = claims.Email
	}

	updated, err := h.service.UpdateUser(r.Context(), u)
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}

	res := toUserRes(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(r.Context(), i); err != nil {
		writeServiceError(w, err, "user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if errors.Is(err, model.ErrInvalidCredentials) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		writeServiceError(w, err, "user")
		return
	}

	res := dto.LoginUserRes{
		SessionID:             result.SessionID,
		AccessToken:           result.AccessToken,
		RefreshToken:          result.RefreshToken,
		AccessTokenExpiresAt:  result.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: result.RefreshTokenExpiresAt,
		User:                  toUserRes(result.User),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := claimsFromContext(r.Context())
	if !ok || claims.SessionID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteSession(r.Context(), claims.SessionID); err != nil {
		writeServiceError(w, err, "session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RenewAccessTokenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	accessToken, expiresAt, err := h.service.RenewAccessToken(r.Context(), req.RefreshToken)
	if errors.Is(err, model.ErrInvalidCredentials) || errors.Is(err, model.ErrNotFound) {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	if err != nil {
		writeServiceError(w, err, "session")
		return
	}

	res := dto.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "session id is required", http.StatusBadRequest)
		return
	}

	if !h.ownsSession(w, r, id) {
		return
	}

	if err := h.service.RevokeSession(r.Context(), id); err != nil {
		writeServiceError(w, err, "session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ownsSession writes an error response and returns false unless the caller owns the session (or is an admin)
func (h *Handler) ownsSession(w http.ResponseWriter, r *http.Request, id string) bool {
	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}

	session, err := h.service.GetSession(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "session")
		return false
	}

	if !claims.IsAdmin && session.UserEmail != claims.Email {
		http.Error(w, "session not found", http.StatusNotFound)
		return false
	}

	return true
}
