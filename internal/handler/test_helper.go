package handler

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
)

type fakeService struct {
	product *model.Product
	products []*model.Product
	order *model.Order
	orders []*model.Order
	err error
}


func (f *fakeService) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	return f.product, f.err
}

func (f *fakeService) CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	return f.product, f.err
}

func (f *fakeService) ListProducts(ctx context.Context) ([]*model.Product, error) {
	return f.products, f.err
}

func (f *fakeService) UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	return f.product, f.err
}

func (f *fakeService) DeleteProduct(ctx context.Context, id int64) error {
	return f.err
}

func (f *fakeService) CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error) {
	return f.order, f.err
}

func (f *fakeService) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	return f.order, f.err
}

func (f *fakeService) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return f.orders, f.err
}

func (f *fakeService) DeleteOrder(ctx context.Context, id int64) error {
	return f.err
}