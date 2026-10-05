package service

import (
	"context"
	"time"

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
	ListOrdersByUser(ctx context.Context, userID int64) ([]*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, from, to model.OrderStatus) error
	DeleteOrder(ctx context.Context, id int64) error

	// users
	CreateUser(ctx context.Context, u *model.User) (*model.User, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, u *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error

	// notifications
	ClaimNotifications(ctx context.Context, limit int, lease time.Duration) ([]*model.Notification, error)
	MarkNotificationSent(ctx context.Context, id int64) error
	MarkNotificationFailed(ctx context.Context, id int64, reason string, maxAttempts int) error

	// sessions
	CreateSession(ctx context.Context, s *model.Session) (*model.Session, error)
	GetSession(ctx context.Context, id string) (*model.Session, error)
	RevokeSession(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, id string) error
}
