package service

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
)

type Storer interface {
	// products
	CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error)
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	ListProducts(ctx context.Context) ([]*model.Product, error)
	UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, id int64) error

	// orders
	CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	ListOrders(ctx context.Context) ([]*model.Order, error)
	DeleteOrder(ctx context.Context, id int64) error

	// users
	CreateUser(ctx context.Context, u *model.User) (*model.User, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, u *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error

	// sessions
	CreateSession(ctx context.Context, s *model.Session) (*model.Session, error)
	GetSession(ctx context.Context, id string) (*model.Session, error)
	RevokeSession(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, id string) error
}
