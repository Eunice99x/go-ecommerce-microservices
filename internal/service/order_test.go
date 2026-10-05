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
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
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

				s := &Service{
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
			test: func(t *testing.T) {
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

				s := &Service{
					storer: &fakeStore,
				}

				_, err := s.CreateOrder(t.Context(), order)

				require.Error(t, err)
			},
		},
		{
			name: "invalid quantity",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:           1,
					Name:         "Iphone",
					Price:        999,
					CountInStock: 10,
				}

				order := &model.Order{
					PaymentMethod: "cash",
					Items: []model.OrderItem{
						{
							ProductID: p.ID,
							Quantity:  0,
						},
					},
				}

				fakeStore := fakeStorer{
					product: p,
				}

				s := &Service{
					storer: &fakeStore,
				}

				_, err := s.CreateOrder(t.Context(), order)

				require.Error(t, err)
			},
		},
		{
			name: "not enough stock",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:           1,
					Name:         "Iphone",
					Price:        999,
					CountInStock: 1,
				}

				order := &model.Order{
					PaymentMethod: "cash",
					Items: []model.OrderItem{
						{
							ProductID: p.ID,
							Quantity:  5,
						},
					},
				}

				fakeStore := fakeStorer{
					product: p,
				}

				s := &Service{
					storer: &fakeStore,
				}

				_, err := s.CreateOrder(t.Context(), order)

				require.Error(t, err)
			},
		},
		{
			name: "failed creating order",
			test: func(t *testing.T) {
				order := &model.Order{
					PaymentMethod: "cash",
				}

				fakeStore := fakeStorer{
					err: fmt.Errorf("error creating order"),
				}

				s := &Service{
					storer: &fakeStore,
				}

				_, err := s.CreateOrder(t.Context(), order)

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
					TotalPrice:    1011.1,
				}

				s := &Service{
					storer: &fakeStorer{order: order},
				}

				got, err := s.GetOrder(t.Context(), 1)

				require.NoError(t, err)
				require.Equal(t, int64(1), got.ID)
				require.Equal(t, "cash", got.PaymentMethod)
			},
		},
		{
			name: "failed getting order",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error getting order")},
				}

				_, err := s.GetOrder(t.Context(), 1)

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

func TestListOrders(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				orders := []*model.Order{
					{ID: 1, PaymentMethod: "cash"},
					{ID: 2, PaymentMethod: "card"},
				}

				s := &Service{
					storer: &fakeStorer{orders: orders},
				}

				got, err := s.ListOrders(t.Context())

				require.NoError(t, err)
				require.Len(t, got, 2)
				require.Equal(t, int64(1), got[0].ID)
				require.Equal(t, int64(2), got[1].ID)
			},
		},
		{
			name: "failed listing orders",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error listing orders")},
				}

				_, err := s.ListOrders(t.Context())

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

func TestListOrdersByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := &Service{
			storer: &fakeStorer{orders: []*model.Order{{ID: 1, UserID: 7}}},
		}

		got, err := s.ListOrdersByUser(t.Context(), 7)

		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("failed listing orders", func(t *testing.T) {
		s := &Service{
			storer: &fakeStorer{err: fmt.Errorf("error listing orders")},
		}

		_, err := s.ListOrdersByUser(t.Context(), 7)

		require.Error(t, err)
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
				s := &Service{
					storer: &fakeStorer{},
				}

				err := s.DeleteOrder(t.Context(), 1)

				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error deleting order")},
				}

				err := s.DeleteOrder(t.Context(), 1)

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

func TestUpdateOrderStatus(t *testing.T) {
	tcs := []struct {
		name    string
		current model.OrderStatus
		next    model.OrderStatus
		wantErr error
	}{
		{name: "pending to shipped", current: model.OrderPending, next: model.OrderShipped},
		{name: "shipped to delivered", current: model.OrderShipped, next: model.OrderDelivered},
		{name: "skipping a step", current: model.OrderPending, next: model.OrderDelivered, wantErr: model.ErrInvalidStatusTransition},
		{name: "going backwards", current: model.OrderShipped, next: model.OrderPending, wantErr: model.ErrInvalidStatusTransition},
		{name: "same status", current: model.OrderShipped, next: model.OrderShipped, wantErr: model.ErrInvalidStatusTransition},
		{name: "already delivered", current: model.OrderDelivered, next: model.OrderShipped, wantErr: model.ErrInvalidStatusTransition},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s := &Service{
				storer: &fakeStorer{order: &model.Order{ID: 1, Status: tc.current}},
			}

			got, err := s.UpdateOrderStatus(t.Context(), 1, tc.next)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.next, got.Status)
		})
	}

	t.Run("failed getting order", func(t *testing.T) {
		s := &Service{
			storer: &fakeStorer{err: fmt.Errorf("error getting order")},
		}

		_, err := s.UpdateOrderStatus(t.Context(), 1, model.OrderShipped)

		require.Error(t, err)
	})
}

func TestCreateOrderInvalidArgument(t *testing.T) {
	tcs := []struct {
		name   string
		storer *fakeStorer
		item   model.OrderItem
	}{
		{name: "unknown product", storer: &fakeStorer{err: fmt.Errorf("error getting product: %w", model.ErrNotFound)}, item: model.OrderItem{ProductID: 9, Quantity: 1}},
		{name: "zero quantity", storer: &fakeStorer{product: &model.Product{ID: 1, CountInStock: 5}}, item: model.OrderItem{ProductID: 1, Quantity: 0}},
		{name: "not enough stock", storer: &fakeStorer{product: &model.Product{ID: 1, CountInStock: 5}}, item: model.OrderItem{ProductID: 1, Quantity: 6}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s := &Service{storer: tc.storer}

			_, err := s.CreateOrder(t.Context(), &model.Order{Items: []model.OrderItem{tc.item}})

			require.ErrorIs(t, err, model.ErrInvalidArgument)
		})
	}
}
