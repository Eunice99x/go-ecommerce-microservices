package model

import "time"

type Order struct {
	ID            int64       `json:"id"`
	PaymentMethod string      `json:"payment_method" db:"payment_method"`
	TaxPrice      float64     `json:"tax_price" db:"tax_price"`
	ShippingPrice float64     `json:"shipping_price" db:"shipping_price"`
	TotalPrice    float64     `json:"total_price" db:"total_price"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at" db:"updated_at"`
	Items         []OrderItem `json:"items"`
}
