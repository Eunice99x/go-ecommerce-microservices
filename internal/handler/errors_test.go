package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestWriteServiceError(t *testing.T) {
	tcs := []struct {
		name     string
		err      error
		wantCode int
		wantBody string
	}{
		{name: "not found", err: fmt.Errorf("error getting order: %w", model.ErrNotFound), wantCode: http.StatusNotFound, wantBody: "order not found"},
		{name: "already exists", err: model.ErrAlreadyExists, wantCode: http.StatusConflict, wantBody: "order already exists"},
		{name: "invalid argument", err: fmt.Errorf("%w: not enough stock", model.ErrInvalidArgument), wantCode: http.StatusBadRequest, wantBody: "invalid argument: not enough stock"},
		{name: "invalid transition", err: fmt.Errorf("%w: shipped -> pending", model.ErrInvalidStatusTransition), wantCode: http.StatusConflict},
		{name: "invalid credentials", err: model.ErrInvalidCredentials, wantCode: http.StatusUnauthorized},
		{name: "internal error is hidden", err: errors.New(`pq: relation "orders" does not exist`), wantCode: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeServiceError(rec, tc.err, "order")

			require.Equal(t, tc.wantCode, rec.Code)
			require.NotContains(t, rec.Body.String(), "pq:")
			if tc.wantBody != "" {
				require.Equal(t, tc.wantBody+"\n", rec.Body.String())
			}
		})
	}
}

// handler-level checks that domain errors from the service reach the client as the right status
func TestDomainErrorStatuses(t *testing.T) {
	withID := func(req *http.Request, id string) *http.Request {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)

		return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	t.Run("missing product is 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h := &Handler{service: &fakeService{err: fmt.Errorf("error getting product: %w", model.ErrNotFound)}}

		h.GetProduct(rec, withID(httptest.NewRequest(http.MethodGet, "/products/9", nil), "9"))

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("missing order is 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h := &Handler{service: &fakeService{err: fmt.Errorf("error getting order: %w", model.ErrNotFound)}}

		req := withClaims(withID(httptest.NewRequest(http.MethodGet, "/orders/9", nil), "9"), testClaims)
		h.GetOrder(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("duplicate email is 409", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h := &Handler{service: &fakeService{err: fmt.Errorf("error creating user: %w", model.ErrAlreadyExists)}}

		body := bytes.NewBufferString(`{"name":"y","email":"taken@example.com","password":"secret123"}`)
		h.CreateUser(rec, httptest.NewRequest(http.MethodPost, "/users", body))

		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("order for unknown product is 400", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h := &Handler{service: &fakeService{err: fmt.Errorf("%w: product 9 does not exist", model.ErrInvalidArgument)}}

		body := bytes.NewBufferString(`{"payment_method":"card","items":[{"quantity":1,"product_id":9}]}`)
		h.CreateOrder(rec, withClaims(httptest.NewRequest(http.MethodPost, "/orders", body), testClaims))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "product 9 does not exist")
	})

	t.Run("login during an outage is 500, not 401", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h := &Handler{service: &fakeService{err: errors.New("grpc unavailable")}}

		body := bytes.NewBufferString(`{"email":"y@example.com","password":"secret123"}`)
		h.LoginUser(rec, httptest.NewRequest(http.MethodPost, "/login", body))

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
