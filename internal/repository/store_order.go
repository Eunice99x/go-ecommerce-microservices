package repository

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/jmoiron/sqlx"
)

func (ps *PostgresStorer) CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error) {

	err := ps.execTx(ctx, func(tx *sqlx.Tx) error {
		order, err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("error creating order: %w", err)
		}

		for _, oi := range o.Items {
			oi.OrderID = order.ID

			err = createOrderItem(ctx, tx, oi)
			if err != nil {
				return fmt.Errorf("error creating order item: %w", err)
			}
		}
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	return o, nil
}

func createOrder(ctx context.Context, tx *sqlx.Tx, o *model.Order) (*model.Order, error) {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (:payment_method, :tax_price, :shipping_price, :total_price)", o)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	o.ID = id

	return o, nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi model.OrderItem) (error) {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (:name, :quantity, :image, :price, :product_id, :order_id)", oi)
	if err != nil {
		return fmt.Errorf("error inserting order item: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	oi.ID = id

	return  nil
}

func (ps *PostgresStorer) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var o model.Order
	err := ps.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	var items []model.OrderItem
	err = ps.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %w", err)
	}
	o.Items = items

	return &o, nil
}

func (ps *PostgresStorer) ListOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := ps.db.SelectContext(ctx, &orders, "SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}
	
	for i := range orders {
		var items []model.OrderItem 
		err = ps.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting order items: %w", err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

// Update order status (later)

func (ps *PostgresStorer) DeleteOrder(ctx context.Context, id int64) error {
	err := ps.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order items: %w", err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	return nil
}