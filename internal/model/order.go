package model

import "time"

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderShipped   OrderStatus = "shipped"
	OrderDelivered OrderStatus = "delivered"
)

func (s OrderStatus) Valid() bool {
	switch s {
	case OrderPending, OrderShipped, OrderDelivered:
		return true
	}

	return false
}

type Order struct {
	ID            int64       `json:"id"`
	PaymentMethod string      `json:"payment_method" db:"payment_method"`
	TaxPrice      float64     `json:"tax_price" db:"tax_price"`
	ShippingPrice float64     `json:"shipping_price" db:"shipping_price"`
	TotalPrice    float64     `json:"total_price" db:"total_price"`
	UserID        int64       `json:"user_id" db:"user_id"`
	Status        OrderStatus `json:"status" db:"status"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at" db:"updated_at"`
	Items         []OrderItem `json:"items"`
}
