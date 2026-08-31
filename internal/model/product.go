package model

import "time"

type Product struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Image        string     `json:"image"`
	Category     string     `json:"category"`
	Description  string     `json:"description"`
	Rating       int64      `json:"rating"`
	NumReviews   int64      `json:"num_reviews" db:"num_reviews"`
	Price        float64    `json:"price"`
	CountInStock int64      `json:"count_in_stock" db:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}
