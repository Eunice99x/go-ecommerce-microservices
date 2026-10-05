// Package mapper converts between models and protobuf messages, used by both server and client.
package mapper

import (
	"time"

	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// products

func ProductToReq(p *model.Product) *pb.ProductReq {
	return &pb.ProductReq{
		Id:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func ProductFromReq(p *pb.ProductReq) *model.Product {
	return &model.Product{
		ID:           p.GetId(),
		Name:         p.GetName(),
		Image:        p.GetImage(),
		Category:     p.GetCategory(),
		Description:  p.GetDescription(),
		Rating:       p.GetRating(),
		NumReviews:   p.GetNumReviews(),
		Price:        p.GetPrice(),
		CountInStock: p.GetCountInStock(),
	}
}

func ProductToRes(p *model.Product) *pb.ProductRes {
	return &pb.ProductRes{
		Id:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    toTimestamp(p.CreatedAt),
		UpdatedAt:    toTimestampPtr(p.UpdatedAt),
	}
}

func ProductFromRes(p *pb.ProductRes) *model.Product {
	return &model.Product{
		ID:           p.GetId(),
		Name:         p.GetName(),
		Image:        p.GetImage(),
		Category:     p.GetCategory(),
		Description:  p.GetDescription(),
		Rating:       p.GetRating(),
		NumReviews:   p.GetNumReviews(),
		Price:        p.GetPrice(),
		CountInStock: p.GetCountInStock(),
		CreatedAt:    fromTimestamp(p.GetCreatedAt()),
		UpdatedAt:    fromTimestampPtr(p.GetUpdatedAt()),
	}
}

// orders

func OrderToReq(o *model.Order) *pb.OrderReq {
	return &pb.OrderReq{
		Id:            o.ID,
		Items:         orderItemsToPB(o.Items),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		UserId:        o.UserID,
	}
}

func OrderFromReq(o *pb.OrderReq) *model.Order {
	return &model.Order{
		ID:            o.GetId(),
		Items:         orderItemsFromPB(o.GetItems()),
		PaymentMethod: o.GetPaymentMethod(),
		TaxPrice:      o.GetTaxPrice(),
		ShippingPrice: o.GetShippingPrice(),
		TotalPrice:    o.GetTotalPrice(),
		UserID:        o.GetUserId(),
	}
}

func OrderToRes(o *model.Order) *pb.OrderRes {
	return &pb.OrderRes{
		Id:            o.ID,
		Items:         orderItemsToPB(o.Items),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		UserId:        o.UserID,
		Status:        OrderStatusToPB(o.Status),
		CreatedAt:     toTimestamp(o.CreatedAt),
		UpdatedAt:     toTimestampPtr(o.UpdatedAt),
	}
}

func OrderFromRes(o *pb.OrderRes) *model.Order {
	return &model.Order{
		ID:            o.GetId(),
		Items:         orderItemsFromPB(o.GetItems()),
		PaymentMethod: o.GetPaymentMethod(),
		TaxPrice:      o.GetTaxPrice(),
		ShippingPrice: o.GetShippingPrice(),
		TotalPrice:    o.GetTotalPrice(),
		UserID:        o.GetUserId(),
		Status:        OrderStatusFromPB(o.GetStatus()),
		CreatedAt:     fromTimestamp(o.GetCreatedAt()),
		UpdatedAt:     fromTimestampPtr(o.GetUpdatedAt()),
	}
}

func OrderStatusToPB(s model.OrderStatus) pb.OrderStatus {
	switch s {
	case model.OrderPending:
		return pb.OrderStatus_ORDER_STATUS_PENDING
	case model.OrderShipped:
		return pb.OrderStatus_ORDER_STATUS_SHIPPED
	case model.OrderDelivered:
		return pb.OrderStatus_ORDER_STATUS_DELIVERED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

// OrderStatusFromPB returns "" for UNSPECIFIED, which model.OrderStatus.Valid rejects
func OrderStatusFromPB(s pb.OrderStatus) model.OrderStatus {
	switch s {
	case pb.OrderStatus_ORDER_STATUS_PENDING:
		return model.OrderPending
	case pb.OrderStatus_ORDER_STATUS_SHIPPED:
		return model.OrderShipped
	case pb.OrderStatus_ORDER_STATUS_DELIVERED:
		return model.OrderDelivered
	default:
		return ""
	}
}

func orderItemsToPB(items []model.OrderItem) []*pb.OrderItem {
	res := make([]*pb.OrderItem, 0, len(items))
	for _, i := range items {
		res = append(res, &pb.OrderItem{
			Id:        i.ID,
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductId: i.ProductID,
			OrderId:   i.OrderID,
		})
	}

	return res
}

func orderItemsFromPB(items []*pb.OrderItem) []model.OrderItem {
	res := make([]model.OrderItem, 0, len(items))
	for _, i := range items {
		res = append(res, model.OrderItem{
			ID:        i.GetId(),
			Name:      i.GetName(),
			Quantity:  i.GetQuantity(),
			Image:     i.GetImage(),
			Price:     i.GetPrice(),
			ProductID: i.GetProductId(),
			OrderID:   i.GetOrderId(),
		})
	}

	return res
}

// users

func UserToReq(u *model.User) *pb.UserReq {
	return &pb.UserReq{
		Id:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func UserFromReq(u *pb.UserReq) *model.User {
	return &model.User{
		ID:       u.GetId(),
		Name:     u.GetName(),
		Email:    u.GetEmail(),
		Password: u.GetPassword(),
		IsAdmin:  u.GetIsAdmin(),
	}
}

// UserToRes never includes the password hash
func UserToRes(u *model.User) *pb.UserRes {
	return &pb.UserRes{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: toTimestamp(u.CreatedAt),
		UpdatedAt: toTimestampPtr(u.UpdatedAt),
	}
}

func UserFromRes(u *pb.UserRes) *model.User {
	return &model.User{
		ID:        u.GetId(),
		Name:      u.GetName(),
		Email:     u.GetEmail(),
		IsAdmin:   u.GetIsAdmin(),
		CreatedAt: fromTimestamp(u.GetCreatedAt()),
		UpdatedAt: fromTimestampPtr(u.GetUpdatedAt()),
	}
}

// sessions

func SessionToReq(s *model.Session) *pb.SessionReq {
	return &pb.SessionReq{
		Id:           s.ID,
		UserEmail:    s.UserEmail,
		RefreshToken: s.RefreshToken,
		IsRevoked:    s.IsRevoked,
		ExpiresAt:    toTimestamp(s.ExpiresAt),
	}
}

func SessionFromReq(s *pb.SessionReq) *model.Session {
	return &model.Session{
		ID:           s.GetId(),
		UserEmail:    s.GetUserEmail(),
		RefreshToken: s.GetRefreshToken(),
		IsRevoked:    s.GetIsRevoked(),
		ExpiresAt:    fromTimestamp(s.GetExpiresAt()),
	}
}

func SessionToRes(s *model.Session) *pb.SessionRes {
	return &pb.SessionRes{
		Id:           s.ID,
		UserEmail:    s.UserEmail,
		RefreshToken: s.RefreshToken,
		IsRevoked:    s.IsRevoked,
		CreatedAt:    toTimestamp(s.CreatedAt),
		ExpiresAt:    toTimestamp(s.ExpiresAt),
	}
}

func SessionFromRes(s *pb.SessionRes) *model.Session {
	return &model.Session{
		ID:           s.GetId(),
		UserEmail:    s.GetUserEmail(),
		RefreshToken: s.GetRefreshToken(),
		IsRevoked:    s.GetIsRevoked(),
		CreatedAt:    fromTimestamp(s.GetCreatedAt()),
		ExpiresAt:    fromTimestamp(s.GetExpiresAt()),
	}
}

// notifications

func NotificationToPB(n *model.Notification) *pb.Notification {
	return &pb.Notification{
		Id:          n.ID,
		OrderId:     n.OrderID,
		OrderStatus: OrderStatusToPB(n.OrderStatus),
		UserEmail:   n.UserEmail,
		Attempts:    int32(n.Attempts),
	}
}

func NotificationFromPB(n *pb.Notification) *model.Notification {
	return &model.Notification{
		ID:          n.GetId(),
		OrderID:     n.GetOrderId(),
		OrderStatus: OrderStatusFromPB(n.GetOrderStatus()),
		UserEmail:   n.GetUserEmail(),
		Attempts:    int(n.GetAttempts()),
	}
}

// timestamps: zero values stay unset, otherwise they come back as 1970-01-01

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}

	return timestamppb.New(t)
}

func toTimestampPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}

	return ts.AsTime()
}

func fromTimestampPtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()

	return &t
}
