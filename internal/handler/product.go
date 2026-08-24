package handler

import (
	"encoding/json"
	"net/http"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
)

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p dto.ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), &model.Product{
		Name: p.Name,
		Image: p.Image,
		Category: p.Category,
		Description: p.Description,
		Price: p.Price,
		CountInStock: p.CountInStock,	
	})
	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := dto.ProductRes{
		ID: product.ID,     
		Name: product.Name,
		Image: product.Image,
		Category: product.Category,
		Description: product.Description,
		Rating: product.Rating,
		NumReviews: product.NumReviews,
		Price: product.Price,
		CountInStock: product.CountInStock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}	