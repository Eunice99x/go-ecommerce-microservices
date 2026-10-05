package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				order := &model.Order{
					ID:            1,
					PaymentMethod: "cash",
					TaxPrice:      10,
					ShippingPrice: 20,
					TotalPrice:    1029,
					Items: []model.OrderItem{
						{
							ID:        1,
							Name:      "Iphone",
							Quantity:  1,
							Image:     "https://example.com",
							Price:     999,
							ProductID: 1,
							OrderID:   1,
						},
					},
				}

				body := []byte(`{
					"payment_method": "cash",
					"items": [
						{
							"quantity": 1,
							"product_id": 1
						}
					]
				}`)

				req := httptest.NewRequest(
					http.MethodPost,
					"/orders",
					bytes.NewReader(body),
				)

				req.Header.Set("Content-Type", "application/json")

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					order: order,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateOrder(rec, req)

				require.Equal(t, http.StatusCreated, rec.Code)

				var res dto.OrderRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, int64(1), res.ID)
				require.Equal(t, "cash", res.PaymentMethod)
				require.Equal(t, float64(1029), res.TotalPrice)
				require.Len(t, res.Items, 1)
				require.Equal(t, "Iphone", res.Items[0].Name)
			},
		},
		{
			name: "invalid request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/orders",
					bytes.NewReader([]byte(`{"items":`)),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.CreateOrder(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed creating order",
			test: func(t *testing.T) {
				body := []byte(`{
					"payment_method": "cash",
					"items": [
						{
							"quantity": 1,
							"product_id": 1
						}
					]
				}`)

				req := httptest.NewRequest(
					http.MethodPost,
					"/orders",
					bytes.NewReader(body),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error creating order"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateOrder(rec, req)

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

func TestGetOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				order := &model.Order{
					ID:            1,
					UserID:        testClaims.ID,
					PaymentMethod: "cash",
					TotalPrice:    999,
					Items: []model.OrderItem{
						{
							ID:        1,
							Name:      "Iphone",
							Quantity:  1,
							Price:     999,
							ProductID: 1,
							OrderID:   1,
						},
					},
				}

				req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					order: order,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.GetOrder(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res dto.OrderRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, int64(1), res.ID)
				require.Len(t, res.Items, 1)
			},
		},
		{
			name: "order belongs to another user",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{order: &model.Order{ID: 1, UserID: 2}},
				}

				h.GetOrder(rec, req)

				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "invalid order id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.GetOrder(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed getting order",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error getting order"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.GetOrder(rec, req)

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

func TestListOrders(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				orders := []*model.Order{
					{
						ID:            1,
						PaymentMethod: "cash",
						TotalPrice:    999,
					},
					{
						ID:            2,
						PaymentMethod: "card",
						TotalPrice:    1500,
					},
				}

				req := httptest.NewRequest(http.MethodGet, "/orders", nil)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					orders: orders,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.ListOrders(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res []dto.OrderRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Len(t, res, 2)
				require.Equal(t, int64(1), res[0].ID)
				require.Equal(t, int64(2), res[1].ID)
			},
		},
		{
			name: "failed listing orders",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/orders", nil)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error listing orders"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.ListOrders(rec, req)

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

func TestListMyOrders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/me", nil)
		req = withClaims(req, testClaims)
		rec := httptest.NewRecorder()

		h := &Handler{
			service: &fakeService{orders: []*model.Order{{ID: 1, UserID: testClaims.ID}}},
		}

		h.ListMyOrders(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var res []dto.OrderRes
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&res))
		require.Len(t, res, 1)
	})

	t.Run("no orders returns empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/me", nil)
		req = withClaims(req, testClaims)
		rec := httptest.NewRecorder()

		h := &Handler{service: &fakeService{}}

		h.ListMyOrders(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, "[]", rec.Body.String())
	})

	t.Run("unauthenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/me", nil)
		rec := httptest.NewRecorder()

		h := &Handler{service: &fakeService{}}

		h.ListMyOrders(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("failed listing orders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/me", nil)
		req = withClaims(req, testClaims)
		rec := httptest.NewRecorder()

		h := &Handler{service: &fakeService{err: fmt.Errorf("error listing orders")}}

		h.ListMyOrders(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDeleteOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{order: &model.Order{ID: 1, UserID: testClaims.ID}},
				}

				h.DeleteOrder(rec, req)

				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name: "invalid order id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/orders/abc", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.DeleteOrder(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed deleting order",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				req = withClaims(req, testClaims)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					order:     &model.Order{ID: 1, UserID: testClaims.ID},
					deleteErr: fmt.Errorf("error deleting order"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.DeleteOrder(rec, req)

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

func TestUpdateOrderStatus(t *testing.T) {
	newReq := func(id, body string) *http.Request {
		req := httptest.NewRequest(http.MethodPatch, "/orders/"+id+"/status", bytes.NewBufferString(body))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)

		return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	tcs := []struct {
		name     string
		id       string
		body     string
		service  *fakeService
		wantCode int
	}{
		{
			name:     "success",
			id:       "1",
			body:     `{"status": "shipped"}`,
			service:  &fakeService{order: &model.Order{ID: 1, Status: model.OrderShipped}},
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid order id",
			id:       "abc",
			body:     `{"status": "shipped"}`,
			service:  &fakeService{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid body",
			id:       "1",
			body:     `{`,
			service:  &fakeService{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "unknown status",
			id:       "1",
			body:     `{"status": "lost"}`,
			service:  &fakeService{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid transition",
			id:       "1",
			body:     `{"status": "pending"}`,
			service:  &fakeService{updateErr: fmt.Errorf("%w: shipped -> pending", model.ErrInvalidStatusTransition)},
			wantCode: http.StatusConflict,
		},
		{
			name:     "service error",
			id:       "1",
			body:     `{"status": "shipped"}`,
			service:  &fakeService{updateErr: fmt.Errorf("error updating order status")},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h := &Handler{service: tc.service}

			h.UpdateOrderStatus(rec, newReq(tc.id, tc.body))

			require.Equal(t, tc.wantCode, rec.Code)

			if tc.wantCode == http.StatusOK {
				var res dto.OrderRes
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&res))
				require.Equal(t, "shipped", res.Status)
			}
		})
	}
}
