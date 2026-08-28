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



func TestCreateProduct(t *testing.T) {
	tcs := []struct{
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				p := &model.Product{
					ID: 1,
					Name:"Iphone",
					Image: "https://exapmle.com",
					Category: "Electronics",
					Rating: 4,
					NumReviews: 14,
					Price: 999,
					CountInStock: 1234,
				}

				body, err := json.Marshal(p)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					product: p,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateProduct(rec, req)
				require.Equal(t, http.StatusCreated, rec.Code)

				var res dto.ProductRes

				err = json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, int64(1), res.ID)
				require.Equal(t, "Iphone", res.Name)
				require.Equal(t, float64(999), res.Price)

			},
		},
		{
			name: "failed decoding request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost,
					"/products",
					bytes.NewBufferString(`{"name":`),
				)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateProduct(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed creating product",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:           1,
					Name:         "Iphone",
					Image:        "https://example.com",
					Category:     "Electronics",
					Price:        999,
					CountInStock: 1234,
				}

				body, err := json.Marshal(p)
				require.NoError(t, err)

				req := httptest.NewRequest(
					http.MethodPost,
					"/products",
					bytes.NewReader(body),
				)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error creating product"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.CreateProduct(rec, req)

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

func TestGetProduct(t *testing.T){
	tcs := []struct{
		name string
		test func(*testing.T)
	}{
		{	
			name: "success",
			test: func(t *testing.T) {

				p := &model.Product{
					ID: 1,
					Name:"Iphone",
					Image: "https://exapmle.com",
					Category: "Electronics",
					Rating: 4,
					NumReviews: 14,
					Price: 999,
					CountInStock: 1234,
				}

				req := httptest.NewRequest(http.MethodGet, "/products/{1}", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					product: p,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.GetProduct(rec, req)
				require.Equal(t, http.StatusOK, rec.Code)
			},

		},
		{
			name: "invalid product id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{}

				h := &Handler{
					service: &fakeS,
				}

				h.GetProduct(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/products/1", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error getting product"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.GetProduct(rec, req)

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

func TestListProducts(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				products := []*model.Product{
					{
						ID:           1,
						Name:         "Iphone",
						Image:        "https://example.com",
						Category:     "Electronics",
						Price:        999,
						CountInStock: 1234,
					},
					{
						ID:           2,
						Name:         "Macbook",
						Image:        "https://example.com",
						Category:     "Electronics",
						Price:        1500,
						CountInStock: 100,
					},
				}

				req := httptest.NewRequest(http.MethodGet, "/products", nil)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					products: products,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.ListProducts(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res []dto.ProductRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Len(t, res, 2)
				require.Equal(t, "Iphone", res[0].Name)
				require.Equal(t, "Macbook", res[1].Name)
			},
		},
		{
			name: "failed listing products",
			test: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/products", nil)
				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error listing products"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.ListProducts(rec, req)

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

func TestUpdateProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:           1,
					Name:         "Iphone",
					Image:        "https://example.com",
					Category:     "Electronics",
					Description:  "old description",
					Price:        999,
					CountInStock: 1234,
				}

				body := []byte(`{
					"name": "Updated Iphone",
					"price": 1200
				}`)

				req := httptest.NewRequest(
					http.MethodPut,
					"/products/1",
					bytes.NewReader(body),
				)

				req.Header.Set("Content-Type", "application/json")

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					product: p,
				}

				h := &Handler{
					service: &fakeS,
				}

				h.UpdateProduct(rec, req)

				require.Equal(t, http.StatusOK, rec.Code)

				var res dto.ProductRes

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.Equal(t, "Updated Iphone", res.Name)
				require.Equal(t, float64(1200), res.Price)
			},
		},
		{
			name: "invalid product id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPut,
					"/products/abc",
					bytes.NewReader([]byte(`{"name":"Iphone"}`)),
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.UpdateProduct(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "invalid request body",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPut,
					"/products/1",
					bytes.NewReader([]byte(`{"name":`)),
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.UpdateProduct(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPut,
					"/products/1",
					bytes.NewReader([]byte(`{"name":"Updated Iphone"}`)),
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error getting product"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.UpdateProduct(rec, req)

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

func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/products/1",
					nil,
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.DeleteProduct(rec, req)

				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name: "invalid product id",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/products/abc",
					nil,
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "abc")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				h := &Handler{
					service: &fakeService{},
				}

				h.DeleteProduct(rec, req)

				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodDelete,
					"/products/1",
					nil,
				)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", "1")

				req = req.WithContext(
					context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
				)

				rec := httptest.NewRecorder()

				fakeS := fakeService{
					err: fmt.Errorf("error deleting product"),
				}

				h := &Handler{
					service: &fakeS,
				}

				h.DeleteProduct(rec, req)

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