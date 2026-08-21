package model

type OrderItem struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Quantity  int64   `json:"quantity"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	ProductID int64   `json:"product_id"`
	OrderID   int64   `json:"order_id"`
}
