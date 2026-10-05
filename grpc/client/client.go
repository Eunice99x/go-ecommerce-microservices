// Package client implements the HTTP handler's Services interface over gRPC.
// The API holds no business logic or DB access; it forwards to the gRPC server.
package client

import (
	"context"
	"time"

	"github.com/eunice99x/goMicro/grpc/mapper"
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/service"
)

type Client struct {
	client pb.EcommClient
}

func New(client pb.EcommClient) *Client {
	return &Client{
		client: client,
	}
}

// products

func (c *Client) CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	res, err := c.client.CreateProduct(ctx, mapper.ProductToReq(p))
	if err != nil {
		return nil, err
	}

	return mapper.ProductFromRes(res), nil
}

func (c *Client) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	res, err := c.client.GetProduct(ctx, &pb.ProductReq{Id: id})
	if err != nil {
		return nil, err
	}

	return mapper.ProductFromRes(res), nil
}

func (c *Client) ListProducts(ctx context.Context) ([]*model.Product, error) {
	res, err := c.client.ListProducts(ctx, &pb.ProductReq{})
	if err != nil {
		return nil, err
	}

	ps := make([]*model.Product, 0, len(res.GetProducts()))
	for _, p := range res.GetProducts() {
		ps = append(ps, mapper.ProductFromRes(p))
	}

	return ps, nil
}

func (c *Client) UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	res, err := c.client.UpdateProduct(ctx, mapper.ProductToReq(p))
	if err != nil {
		return nil, err
	}

	return mapper.ProductFromRes(res), nil
}

func (c *Client) DeleteProduct(ctx context.Context, id int64) error {
	_, err := c.client.DeleteProduct(ctx, &pb.ProductReq{Id: id})

	return err
}

// orders

func (c *Client) CreateOrder(ctx context.Context, o *model.Order) (*model.Order, error) {
	res, err := c.client.CreateOrder(ctx, mapper.OrderToReq(o))
	if err != nil {
		return nil, err
	}

	return mapper.OrderFromRes(res), nil
}

func (c *Client) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	res, err := c.client.GetOrder(ctx, &pb.OrderReq{Id: id})
	if err != nil {
		return nil, err
	}

	return mapper.OrderFromRes(res), nil
}

func (c *Client) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return c.listOrders(ctx, &pb.OrderReq{})
}

func (c *Client) ListOrdersByUser(ctx context.Context, userID int64) ([]*model.Order, error) {
	return c.listOrders(ctx, &pb.OrderReq{UserId: userID})
}

func (c *Client) listOrders(ctx context.Context, req *pb.OrderReq) ([]*model.Order, error) {
	res, err := c.client.ListOrders(ctx, req)
	if err != nil {
		return nil, err
	}

	orders := make([]*model.Order, 0, len(res.GetOrders()))
	for _, o := range res.GetOrders() {
		orders = append(orders, mapper.OrderFromRes(o))
	}

	return orders, nil
}

func (c *Client) DeleteOrder(ctx context.Context, id int64) error {
	_, err := c.client.DeleteOrder(ctx, &pb.OrderReq{Id: id})

	return err
}

// users

func (c *Client) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	res, err := c.client.CreateUser(ctx, mapper.UserToReq(u))
	if err != nil {
		return nil, err
	}

	return mapper.UserFromRes(res), nil
}

func (c *Client) GetUser(ctx context.Context, email string) (*model.User, error) {
	res, err := c.client.GetUser(ctx, &pb.UserReq{Email: email})
	if err != nil {
		return nil, err
	}

	return mapper.UserFromRes(res), nil
}

func (c *Client) ListUsers(ctx context.Context) ([]*model.User, error) {
	res, err := c.client.ListUsers(ctx, &pb.UserReq{})
	if err != nil {
		return nil, err
	}

	us := make([]*model.User, 0, len(res.GetUsers()))
	for _, u := range res.GetUsers() {
		us = append(us, mapper.UserFromRes(u))
	}

	return us, nil
}

func (c *Client) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	res, err := c.client.UpdateUser(ctx, mapper.UserToReq(u))
	if err != nil {
		return nil, err
	}

	return mapper.UserFromRes(res), nil
}

func (c *Client) DeleteUser(ctx context.Context, id int64) error {
	_, err := c.client.DeleteUser(ctx, &pb.UserReq{Id: id})

	return err
}

// auth

func (c *Client) LoginUser(ctx context.Context, email, password string) (*service.LoginResult, error) {
	res, err := c.client.LoginUser(ctx, &pb.LoginUserReq{Email: email, Password: password})
	if err != nil {
		return nil, err
	}

	return &service.LoginResult{
		User:                  mapper.UserFromRes(res.GetUser()),
		SessionID:             res.GetSessionId(),
		AccessToken:           res.GetAccessToken(),
		RefreshToken:          res.GetRefreshToken(),
		AccessTokenExpiresAt:  res.GetAccessTokenExpiresAt().AsTime(),
		RefreshTokenExpiresAt: res.GetRefreshTokenExpiresAt().AsTime(),
	}, nil
}

func (c *Client) RenewAccessToken(ctx context.Context, refreshToken string) (string, time.Time, error) {
	res, err := c.client.RenewAccessToken(ctx, &pb.RenewAccessTokenReq{RefreshToken: refreshToken})
	if err != nil {
		return "", time.Time{}, err
	}

	return res.GetAccessToken(), res.GetAccessTokenExpiresAt().AsTime(), nil
}

// sessions

func (c *Client) GetSession(ctx context.Context, id string) (*model.Session, error) {
	res, err := c.client.GetSession(ctx, &pb.SessionReq{Id: id})
	if err != nil {
		return nil, err
	}

	return mapper.SessionFromRes(res), nil
}

func (c *Client) RevokeSession(ctx context.Context, id string) error {
	_, err := c.client.RevokeSession(ctx, &pb.SessionReq{Id: id})

	return err
}

func (c *Client) DeleteSession(ctx context.Context, id string) error {
	_, err := c.client.DeleteSession(ctx, &pb.SessionReq{Id: id})

	return err
}
