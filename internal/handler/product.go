package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/go-chi/chi/v5"
)

func toProductModel(p dto.ProductReq) *model.Product {
	return &model.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toProductRes(p *model.Product) dto.ProductRes {
	return dto.ProductRes{
		ID:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func patchProductReq(p *model.Product, req dto.ProductReq) {
	if req.Name != "" {
		p.Name = req.Name
	}

	if req.Image != "" {
		p.Image = req.Image
	}

	if req.Category != "" {
		p.Category = req.Category
	}

	if req.Description != "" {
		p.Description = req.Description
	}

	if req.Price != 0 {
		p.Price = req.Price
	}

	if req.CountInStock != 0 {
		p.CountInStock = req.CountInStock
	}

	now := time.Now()

	p.UpdatedAt = &now
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(
		r.Context(),
		toProductModel(req),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating product: %v", err), http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProduct(r.Context(), i)
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		http.Error(w, "error listing products", http.StatusInternalServerError)
		return
	}

	res := make([]dto.ProductRes, 0, len(products))

	for _, product := range products {
		res = append(res, toProductRes(product))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req dto.ProductReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetProduct(r.Context(), i)
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	patchProductReq(p, req)

	updated, err := h.service.UpdateProduct(r.Context(), p)
	if err != nil {
		http.Error(w, "error updating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), i); err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
