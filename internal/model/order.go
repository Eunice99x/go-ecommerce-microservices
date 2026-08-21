package model

import "time"

type Order struct {
	ID            int64       `json:"id"`
	PaymentMethod string      `json:"payment_method"`
	TaxPrice      float64     `json:"tax_price"`
	ShippingPrice float64     `json:"shipping_price"`
	TotalPrice    float64     `json:"total_price"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at"`
	Items         []OrderItem `json:"items"`
}
