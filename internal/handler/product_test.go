package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
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
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}