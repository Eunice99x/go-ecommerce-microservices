package service

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
)

type fakeStorer struct {
	product *model.Product
	order   *model.Order
	err     error
}

func (f *fakeStorer) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	return f.product, f.err
}

func (f *fakeStorer) CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	return f.product, f.err
}

func (f *fakeStorer) ListProducts(ctx context.Context) ([]*model.Product, error) {
	return nil, f.err
}

func (f *fakeStorer) UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	return f.product, f.err
}

func (f *fakeStorer) DeleteProduct(ctx context.Context, id int64) error {
	return f.err
}

func (f *fakeStorer) CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error) {
	return f.order, f.err
}

func (f *fakeStorer) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	return f.order, f.err
}

func (f *fakeStorer) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return nil, f.err
}

func (f *fakeStorer) DeleteOrder(ctx context.Context, id int64) error {
	return f.err
}
