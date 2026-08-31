package service

import (
	"fmt"
	"testing"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *Service, *model.Order)
	}{
		{
			name: "success",
			test: func(t *testing.T, s *Service, o *model.Order) {
				p := &model.Product{
					ID:           1,
					Name:         "Iphone",
					Image:        "https://exapmle.com",
					Category:     "Electronics",
					Rating:       4,
					NumReviews:   14,
					Price:        999,
					CountInStock: 1234,
				}

				order := &model.Order{
					ID:            1,
					PaymentMethod: "cash",
					TaxPrice:      0.1,
					ShippingPrice: 12,
				}

				oi := &model.OrderItem{
					ID:        1,
					Quantity:  1,
					ProductID: p.ID,
					OrderID:   order.ID,
				}

				order.Items = []model.OrderItem{*oi}

				fakeStore := fakeStorer{
					product: p,
					order:   order,
				}

				s = &Service{
					storer: &fakeStore,
				}

				got, err := s.CreateOrder(t.Context(), order)

				require.NoError(t, err)
				require.Equal(t, "Iphone", got.Items[0].Name)
				require.Equal(t, float64(999), got.Items[0].Price)
				require.Equal(t, float64(1011.1), got.TotalPrice)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T, s *Service, o *model.Order) {
				order := &model.Order{
					PaymentMethod: "cash",
					Items: []model.OrderItem{
						{
							ProductID: 1,
							Quantity:  1,
						},
					},
				}

				fakeStore := fakeStorer{
					err: fmt.Errorf("error getting product"),
				}

				s = &Service{
					storer: &fakeStore,
				}

				_, err := s.CreateOrder(t.Context(), order)

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t, nil, nil)
		})
	}
}
