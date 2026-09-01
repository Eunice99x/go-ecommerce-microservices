package handler

import (
	"context"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/service"
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

	CreateUser(ctx context.Context, u *model.User) (*model.User, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, u *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error

	LoginUser(ctx context.Context, email, password string) (*service.LoginResult, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (string, time.Time, error)
	RevokeSession(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, id string) error
}
