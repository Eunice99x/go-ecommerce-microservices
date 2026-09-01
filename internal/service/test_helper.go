package service

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
)

type fakeStorer struct {
	product *model.Product
	order   *model.Order
	user    *model.User
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

// order fake funcs

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

// user fake funcs

func (f *fakeStorer) GetUser(ctx context.Context, email string) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeStorer) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeStorer) ListUsers(ctx context.Context) ([]*model.User, error) {
	return nil, f.err
}

func (f *fakeStorer) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeStorer) DeleteUser(ctx context.Context, id int64) error {
	return f.err
}
