// Package server exposes the service layer over gRPC. Business logic stays in
// internal/service; this package only maps between protobuf and domain models.
package server

import (
	"context"
	"time"

	"github.com/eunice99x/goMicro/grpc/mapper"
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service *service.Service
	pb.UnimplementedEcommServer
}

func NewServer(svc *service.Service) *Server {
	return &Server{
		service: svc,
	}
}

// products

func (s *Server) CreateProduct(ctx context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	p, err := s.service.CreateProduct(ctx, mapper.ProductFromReq(req))
	if err != nil {
		return nil, err
	}

	return mapper.ProductToRes(p), nil
}

func (s *Server) GetProduct(ctx context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	p, err := s.service.GetProduct(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return mapper.ProductToRes(p), nil
}

func (s *Server) ListProducts(ctx context.Context, _ *pb.ProductReq) (*pb.ListProductRes, error) {
	ps, err := s.service.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.ProductRes, 0, len(ps))
	for _, p := range ps {
		res = append(res, mapper.ProductToRes(p))
	}

	return &pb.ListProductRes{Products: res}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	p := mapper.ProductFromReq(req)

	// ProductReq has no updated_at, so the server stamps it
	now := time.Now()
	p.UpdatedAt = &now

	if _, err := s.service.UpdateProduct(ctx, p); err != nil {
		return nil, err
	}

	// re-read so fields the request doesn't carry (created_at) come back filled in
	updated, err := s.service.GetProduct(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	return mapper.ProductToRes(updated), nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	if err := s.service.DeleteProduct(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &pb.ProductRes{}, nil
}

// orders

func (s *Server) CreateOrder(ctx context.Context, req *pb.OrderReq) (*pb.OrderRes, error) {
	o, err := s.service.CreateOrder(ctx, mapper.OrderFromReq(req))
	if err != nil {
		return nil, err
	}

	return mapper.OrderToRes(o), nil
}

func (s *Server) GetOrder(ctx context.Context, req *pb.OrderReq) (*pb.OrderRes, error) {
	o, err := s.service.GetOrder(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return mapper.OrderToRes(o), nil
}

// ListOrders returns every order, or only one user's orders when user_id is set
func (s *Server) ListOrders(ctx context.Context, req *pb.OrderReq) (*pb.ListOrderRes, error) {
	var (
		orders []*model.Order
		err    error
	)

	if req.GetUserId() != 0 {
		orders, err = s.service.ListOrdersByUser(ctx, req.GetUserId())
	} else {
		orders, err = s.service.ListOrders(ctx)
	}
	if err != nil {
		return nil, err
	}

	res := make([]*pb.OrderRes, 0, len(orders))
	for _, o := range orders {
		res = append(res, mapper.OrderToRes(o))
	}

	return &pb.ListOrderRes{Orders: res}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *pb.OrderReq) (*pb.OrderRes, error) {
	if err := s.service.DeleteOrder(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &pb.OrderRes{}, nil
}

// users

func (s *Server) CreateUser(ctx context.Context, req *pb.UserReq) (*pb.UserRes, error) {
	u, err := s.service.CreateUser(ctx, mapper.UserFromReq(req))
	if err != nil {
		return nil, err
	}

	return mapper.UserToRes(u), nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.UserReq) (*pb.UserRes, error) {
	u, err := s.service.GetUser(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return mapper.UserToRes(u), nil
}

func (s *Server) ListUsers(ctx context.Context, _ *pb.UserReq) (*pb.ListUserRes, error) {
	us, err := s.service.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.UserRes, 0, len(us))
	for _, u := range us {
		res = append(res, mapper.UserToRes(u))
	}

	return &pb.ListUserRes{Users: res}, nil
}

// UpdateUser treats an empty password as "unchanged"; the repository keeps the stored hash
func (s *Server) UpdateUser(ctx context.Context, req *pb.UserReq) (*pb.UserRes, error) {
	u := mapper.UserFromReq(req)

	// UserReq has no updated_at, so the server stamps it
	now := time.Now()
	u.UpdatedAt = &now

	if _, err := s.service.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	// re-read so fields the request doesn't carry (created_at) come back filled in
	updated, err := s.service.GetUser(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	return mapper.UserToRes(updated), nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.UserReq) (*pb.UserRes, error) {
	if err := s.service.DeleteUser(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &pb.UserRes{}, nil
}

// auth

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserReq) (*pb.LoginUserRes, error) {
	result, err := s.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginUserRes{
		User:                  mapper.UserToRes(result.User),
		SessionId:             result.SessionID,
		AccessToken:           result.AccessToken,
		RefreshToken:          result.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(result.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(result.RefreshTokenExpiresAt),
	}, nil
}

func (s *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenReq) (*pb.RenewAccessTokenRes, error) {
	accessToken, expiresAt, err := s.service.RenewAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &pb.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(expiresAt),
	}, nil
}

// sessions

func (s *Server) CreateSession(ctx context.Context, req *pb.SessionReq) (*pb.SessionRes, error) {
	sess, err := s.service.CreateSession(ctx, mapper.SessionFromReq(req))
	if err != nil {
		return nil, err
	}

	return mapper.SessionToRes(sess), nil
}

func (s *Server) GetSession(ctx context.Context, req *pb.SessionReq) (*pb.SessionRes, error) {
	sess, err := s.service.GetSession(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return mapper.SessionToRes(sess), nil
}

func (s *Server) RevokeSession(ctx context.Context, req *pb.SessionReq) (*pb.SessionRes, error) {
	if err := s.service.RevokeSession(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &pb.SessionRes{}, nil
}

func (s *Server) DeleteSession(ctx context.Context, req *pb.SessionReq) (*pb.SessionRes, error) {
	if err := s.service.DeleteSession(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &pb.SessionRes{}, nil
}
