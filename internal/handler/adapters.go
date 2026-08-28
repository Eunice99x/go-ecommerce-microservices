package handler

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
)

type Services interface {
	CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error)
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	ListProducts(ctx context.Context) ([]*model.Product, error)
	UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, id int64) error

	CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	ListOrders(ctx context.Context) ([]*model.Order, error)
	DeleteOrder(ctx context.Context, id int64) error
}