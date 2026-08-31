package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

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

				oRows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				oiRows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				s.ExpectQuery(`INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES ($1, $2, $3, $4) RETURNING id, created_at`).WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice).WillReturnRows(oRows)

				s.ExpectQuery(`INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`).WithArgs(oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, oi.OrderID).WillReturnRows(oiRows)

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
		{
			name:"failed to create order",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectBegin()

				s.ExpectQuery(`INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES ($1, $2, $3, $4) RETURNING id, created_at`).WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice).WillReturnError(fmt.Errorf("error inserting order"))

				s.ExpectRollback()

				_, err := ps.CreateOrder(context.Background(), o)
	
				require.Error(t, err)
				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name:"failed to create order item",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectBegin()

				oRows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				s.ExpectQuery(`INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES ($1, $2, $3, $4) RETURNING id, created_at`).WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice).WillReturnRows(oRows)

				s.ExpectQuery(`INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`).WithArgs(oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, oi.OrderID).WillReturnError(fmt.Errorf("error inserting order item"))

				s.ExpectRollback()

				_, err := ps.CreateOrder(context.Background(), o)
	
				require.Error(t, err)
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

func TestGetOrder(t *testing.T) {

	oi := model.OrderItem{
		ID:        1,
		Name:      "iphone",
		Quantity:  1,
		Image:     "www.example.com",
		Price:     1111,
		ProductID: 1,
		OrderID:   1,
	}

	o := &model.Order{
		ID: 1,
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
				
				oiRows := s.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).AddRow(o.ID, oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, oi.OrderID)

				rows := s.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).AddRow(o.ID, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				s.ExpectQuery("SELECT * FROM orders WHERE id=$1").WithArgs(1).WillReturnRows(rows)
				s.ExpectQuery("SELECT * FROM order_items WHERE order_id=$1").WithArgs(1).WillReturnRows(oiRows)

				order, err := ps.GetOrder(context.Background(), 1)


				require.NoError(t, err)
				require.Equal(t, int64(1), order.ID)
				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name:"error getting order",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT * FROM orders WHERE id=$1").WithArgs(1).WillReturnError(fmt.Errorf("error getting order"))

				_, err := ps.GetOrder(context.Background(), 1)

				require.Error(t, err)
				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error getting order item",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock)  {
				rows := s.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).AddRow(o.ID, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)

				s.ExpectQuery("SELECT * FROM orders WHERE id=$1").WithArgs(1).WillReturnRows(rows)
				s.ExpectQuery("SELECT * FROM order_items WHERE order_id=$1").WithArgs(1).WillReturnError(fmt.Errorf("error getting order items"))

				_, err := ps.GetOrder(context.Background(), 1)

				require.Error(t, err)
				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, s sqlmock.Sqlmock)  {
				n := NewPostgresStorer(db)
				tc.test(t, n, s)
			})
		})
	}
}

func TestListOrders(t *testing.T) {

	oi := model.OrderItem{
		ID:        1,
		Name:      "iphone",
		Quantity:  1,
		Image:     "www.example.com",
		Price:     1111,
		ProductID: 1,
		OrderID:   1,
	}

	values := [][]driver.Value{
		{1, "cash", 22.2, 100.0, 122.2, time.Now(), nil},
		{2, "cash", 22.2, 100.0, 122.2, time.Now(), nil},
		{3, "cash", 22.2, 100.0, 122.2, time.Now(), nil},
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock){
				rows := s.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).AddRows(values...)

				s.ExpectQuery("SELECT * FROM orders").WillReturnRows(rows)

				for _, orderID := range []int64{1, 2, 3} {
					oiRows := s.NewRows([]string{"id","name","quantity","image","price","product_id","order_id",}).
					AddRow(oi.ID, oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, orderID)

					s.ExpectQuery("SELECT * FROM order_items WHERE order_id=$1").WithArgs(orderID).WillReturnRows(oiRows)
				}

				orders, err := ps.ListOrders(context.Background())
				require.NoError(t, err)

				require.Len(t, orders, 3)
				require.Equal(t, int64(1), orders[0].ID)
				require.Equal(t, int64(2), orders[1].ID)
				require.Equal(t, int64(3), orders[2].ID)

				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error listing orders",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT * FROM orders").WillReturnError(fmt.Errorf("error listing orders"))

				_, err := ps.ListOrders(context.Background())

				require.Error(t, err)

				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error getting order items",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				now := time.Now()

				rows := s.NewRows([]string{"id","payment_method","tax_price","shipping_price","total_price","created_at","updated_at"}).
				AddRow(1, "cash", 22.2, 100.0, 122.2, now, nil)

				s.ExpectQuery("SELECT * FROM orders").WillReturnRows(rows)
				s.ExpectQuery("SELECT * FROM order_items WHERE order_id=$1").WithArgs(1).WillReturnError(fmt.Errorf("error getting order items"))

				_, err := ps.ListOrders(context.Background())

				require.Error(t, err)

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

func TestDeleteOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectBegin()

				s.ExpectExec("DELETE FROM order_items WHERE order_id=$1").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
				s.ExpectExec("DELETE FROM orders WHERE id=$1").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

				s.ExpectCommit()

				err := ps.DeleteOrder(context.Background(), int64(1))
				require.NoError(t, err)

				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order items",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectBegin()

				s.ExpectExec("DELETE FROM order_items WHERE order_id=$1").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order items"))

				s.ExpectRollback()

				err := ps.DeleteOrder(context.Background(), int64(1))

				require.Error(t, err)

				err = s.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order",
			test: func(t *testing.T, ps *PostgresStorer, s sqlmock.Sqlmock) {
				s.ExpectBegin()

				s.ExpectExec("DELETE FROM order_items WHERE order_id=$1").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

				s.ExpectExec("DELETE FROM orders WHERE id=$1").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order"))

				s.ExpectRollback()

				err := ps.DeleteOrder(context.Background(), int64(1))

				require.Error(t, err)

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