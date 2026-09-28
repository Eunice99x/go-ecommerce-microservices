package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/eunice99x/goMicro/internal/service"
)

type fakeService struct {
	product     *model.Product
	products    []*model.Product
	order       *model.Order
	orders      []*model.Order
	user        *model.User
	session     *model.Session
	users       []*model.User
	loginResult *service.LoginResult
	accessToken string
	expiresAt   time.Time
	err         error

	// set to fail only the update call while the preceding get still succeeds
	updateErr error
	// same idea for delete/revoke calls that are preceded by an ownership lookup
	deleteErr error
	revokeErr error
}

// testClaims is the logged-in (non-admin) user used by handler tests
var testClaims = &auth.Claims{ID: 1, Email: "younes@example.com"}

// withClaims simulates the auth middleware having run
func withClaims(r *http.Request, c *auth.Claims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authKey{}, c))
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
	if f.updateErr != nil {
		return nil, f.updateErr
	}

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

func (f *fakeService) ListOrdersByUser(ctx context.Context, userID int64) ([]*model.Order, error) {
	return f.orders, f.err
}

func (f *fakeService) DeleteOrder(ctx context.Context, id int64) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}

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
	if f.updateErr != nil {
		return nil, f.updateErr
	}

	return f.user, f.err
}

func (f *fakeService) DeleteUser(ctx context.Context, id int64) error {
	return f.err
}

// user login
func (f *fakeService) LoginUser(ctx context.Context, email, password string) (*service.LoginResult, error) {
	return f.loginResult, f.err
}

func (f *fakeService) RenewAccessToken(ctx context.Context, refreshToken string) (string, time.Time, error) {
	return f.accessToken, f.expiresAt, f.err
}

func (f *fakeService) GetSession(ctx context.Context, id string) (*model.Session, error) {
	return f.session, f.err
}

func (f *fakeService) RevokeSession(ctx context.Context, id string) error {
	if f.revokeErr != nil {
		return f.revokeErr
	}

	return f.err
}

func (f *fakeService) DeleteSession(ctx context.Context, id string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}

	return f.err
}
