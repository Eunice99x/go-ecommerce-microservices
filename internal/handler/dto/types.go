package dto

import "time"

type ProductReq struct {
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	CountInStock int64   `json:"count_in_stock"`
}

type ProductRes struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Image        string     `json:"image"`
	Category     string     `json:"category"`
	Description  string     `json:"description"`
	Rating       int64      `json:"rating"`
	NumReviews   int64      `json:"num_reviews"`
	Price        float64    `json:"price"`
	CountInStock int64      `json:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type OrderReq struct {
	Items         []*OrderItemReq `json:"items"`
	PaymentMethod string          `json:"payment_method"`
}

type OrderRes struct {
	ID            int64           `json:"id"`
	Items         []*OrderItemRes `json:"items"`
	PaymentMethod string          `json:"payment_method"`
	TaxPrice      float64        `json:"tax_price"`
	ShippingPrice float64         `json:"shipping_price"`
	TotalPrice    float64         `json:"total_price"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}


type OrderItemReq struct {
	Quantity  int64 `json:"quantity"`
	ProductID int64 `json:"product_id"`
}

type OrderItemRes struct {
	Name      string  `json:"name"`
	Quantity  int64   `json:"quantity"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	ProductID int64   `json:"product_id"`
}