package handler

import (
	"context"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/service"
)

type fakeService struct {
	product     *model.Product
	products    []*model.Product
	order       *model.Order
	orders      []*model.Order
	user        *model.User
	users       []*model.User
	loginResult *service.LoginResult
	accessToken string
	expiresAt   time.Time
	err         error
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

// order fake funcs

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

// user fake funcs

func (f *fakeService) GetUser(ctx context.Context, email string) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeService) CreateUser(ctx context.Context, p *model.User) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return f.users, f.err
}

func (f *fakeService) UpdateUser(ctx context.Context, p *model.User) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeService) DeleteUser(ctx context.Context, id int64) error {
	return f.err
}

// user login
func (f *fakeService) LoginUser(ctx context.Context, email, password string) (*service.LoginResult, error) {
	return nil, f.err
}

func (f *fakeService) RenewAccessToken(ctx context.Context, refreshToken string) (string, time.Time, error) {
	return f.accessToken, f.expiresAt, f.err
}

func (f *fakeService) RevokeSession(ctx context.Context, id string) error {
	return f.err
}

func (f *fakeService) DeleteSession(ctx context.Context, id string) error {
	return f.err
}
