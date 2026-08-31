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
			name: "invalid order id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

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

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
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

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error deleting order"),
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
