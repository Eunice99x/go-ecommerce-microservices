package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
)

func (s *Service) CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error) {
	var totalPrice float64

	for i := range o.Items {
		product, err := s.storer.GetProduct(ctx, o.Items[i].ProductID)
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("%w: product %d does not exist", model.ErrInvalidArgument, o.Items[i].ProductID)
		}
		if err != nil {
			return nil, fmt.Errorf("error getting product %d: %w", o.Items[i].ProductID, err)
		}

		if o.Items[i].Quantity <= 0 {
			return nil, fmt.Errorf("%w: product quantity must be greater than 0", model.ErrInvalidArgument)
		}

		if o.Items[i].Quantity > product.CountInStock {
			return nil, fmt.Errorf("%w: not enough stock for product %d", model.ErrInvalidArgument, product.ID)
		}

		o.Items[i].Name = product.Name
		o.Items[i].Image = product.Image
		o.Items[i].Price = product.Price

		totalPrice += product.Price * float64(o.Items[i].Quantity)
	}

	o.TotalPrice = totalPrice + o.TaxPrice + o.ShippingPrice

	order, err := s.storer.CreateOrder(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	return order, nil
}

func (s *Service) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	return s.storer.GetOrder(ctx, id)
}

func (s *Service) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return s.storer.ListOrders(ctx)
}

func (s *Service) ListOrdersByUser(ctx context.Context, userID int64) ([]*model.Order, error) {
	return s.storer.ListOrdersByUser(ctx, userID)
}

// orders only move forward: pending -> shipped -> delivered
var nextOrderStatus = map[model.OrderStatus]model.OrderStatus{
	model.OrderPending: model.OrderShipped,
	model.OrderShipped: model.OrderDelivered,
}

// UpdateOrderStatus advances an order and queues the customer email in the same transaction
func (s *Service) UpdateOrderStatus(ctx context.Context, id int64, status model.OrderStatus) (*model.Order, error) {
	o, err := s.storer.GetOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	if nextOrderStatus[o.Status] != status {
		return nil, fmt.Errorf("%w: %s -> %s", model.ErrInvalidStatusTransition, o.Status, status)
	}

	if err := s.storer.UpdateOrderStatus(ctx, id, o.Status, status); err != nil {
		return nil, err
	}

	return s.storer.GetOrder(ctx, id)
}

func (s *Service) DeleteOrder(ctx context.Context, id int64) error {
	return s.storer.DeleteOrder(ctx, id)
}
