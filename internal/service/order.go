package service

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
)


func (s *Service) CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error){
	var totalPrice float64

	for i := range o.Items {
		product, err := s.storer.GetProduct(ctx, o.Items[i].ProductID)
		if err != nil {
			return nil, fmt.Errorf(
				"error getting product %d: %w",
				o.Items[i].ProductID,
				err,
			)
		}

		if o.Items[i].Quantity <= 0 {
			return nil, fmt.Errorf("product quantity must be greater than 0")
		}

		if o.Items[i].Quantity > product.CountInStock {
			return nil, fmt.Errorf("not enough stock for product %d", product.ID)
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

func (s *Service) GetOrder(ctx context.Context, id int64) (*model.Order, error){
	return s.storer.GetOrder(ctx, id)
}

func (s *Service) ListOrders(ctx context.Context) ([]*model.Order, error){
	return s.storer.ListOrders(ctx)
}


// will do it later after adding noti
// func (s *Service) UpdateOrder(ctx context.Context, p *model.Order) (*model.Order, error) {
// 	return s.storer.UpdateOrder(ctx, p)
// }

func (s *Service) DeleteOrder(ctx context.Context, id int64) error {
	return s.storer.DeleteOrder(ctx, id)
}
