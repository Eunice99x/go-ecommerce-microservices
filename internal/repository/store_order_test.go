package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {

	p := &model.Product{
		ID:           1,
		Name:         "iphone",
		Image:        "www.example.com",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        100,
		CountInStock: 999,
	}

	oi := model.OrderItem{
		ID:        1,
		Name:      "iphone",
		Quantity:  1,
		Image:     "www.example.com",
		Price:     1111,
		ProductID: p.ID,
		OrderID:   1,
	}

	o := &model.Order{
		PaymentMethod: "cash",
		TaxPrice:      22.2,
		ShippingPrice: 100,
		TotalPrice:    122.2,
		Items:         []model.OrderItem{oi},
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {

				s.ExpectBegin()

				s.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice).WillReturnResult(sqlmock.NewResult(1, 1))

				s.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WithArgs(oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, oi.OrderID).WillReturnResult(sqlmock.NewResult(1, 1))

				s.ExpectCommit()

				co, err := ps.CreateOrder(context.Background(), o)
				if err != nil {
					t.Fatalf("error creating an order: %v", err)
				}

				require.NoError(t, err)
				require.Equal(t, int64(1), co.ID)
				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock) {
				n := NewPostgresStorer(db)
				tc.test(t, n, s)
			})
		})
	}

}
