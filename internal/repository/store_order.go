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

		for i := range o.Items {
			o.Items[i].OrderID = order.ID

			err = createOrderItem(ctx, tx, &o.Items[i])
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
	query := `INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	err := tx.GetContext(ctx, o, query, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.UserID)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %w", err)
	}

	return o, nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi *model.OrderItem) error {

	query := `INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := tx.GetContext(ctx, &oi.ID, query, oi.Name, oi.Quantity, oi.Image, oi.Price, oi.ProductID, oi.OrderID)
	if err != nil {
		return fmt.Errorf("error inserting order item: %w", err)
	}

	return nil
}

func (ps *PostgresStorer) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var o model.Order
	err := ps.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	var items []model.OrderItem
	err = ps.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=$1", o.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %w", err)
	}
	o.Items = items

	return &o, nil
}

func (ps *PostgresStorer) ListOrders(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order
	err := ps.db.SelectContext(ctx, &orders, "SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("error listing orders: %w", err)
	}

	if err := ps.loadOrderItems(ctx, orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (ps *PostgresStorer) ListOrdersByUser(ctx context.Context, userID int64) ([]*model.Order, error) {
	var orders []*model.Order
	err := ps.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("error listing user orders: %w", err)
	}

	if err := ps.loadOrderItems(ctx, orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (ps *PostgresStorer) loadOrderItems(ctx context.Context, orders []*model.Order) error {
	for i := range orders {
		var items []model.OrderItem

		err := ps.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=$1", orders[i].ID)
		if err != nil {
			return fmt.Errorf("error getting order items: %w", err)
		}
		orders[i].Items = items
	}

	return nil
}

// Update order status (later)

func (ps *PostgresStorer) DeleteOrder(ctx context.Context, id int64) error {
	err := ps.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id=$1", id)
		if err != nil {
			return fmt.Errorf("error deleting order items: %w", err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", id)
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
