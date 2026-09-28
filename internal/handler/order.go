package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eunice99x/goMicro/internal/handler/dto"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/go-chi/chi/v5"
)

func toOrderItemModel(oi dto.OrderItemReq) *model.OrderItem {
	return &model.OrderItem{
		Quantity:  oi.Quantity,
		ProductID: oi.ProductID,
	}
}

func toOrderModel(o dto.OrderReq) *model.Order {
	items := make([]model.OrderItem, 0, len(o.Items))

	for _, item := range o.Items {
		items = append(items, *toOrderItemModel(item))
	}

	return &model.Order{
		Items:         items,
		PaymentMethod: o.PaymentMethod,
	}
}

func toOrderItemRes(oi model.OrderItem) dto.OrderItemRes {
	return dto.OrderItemRes{
		Name:      oi.Name,
		Quantity:  oi.Quantity,
		Image:     oi.Image,
		Price:     oi.Price,
		ProductID: oi.ProductID,
	}
}

func toOrderRes(o *model.Order) dto.OrderRes {
	items := make([]dto.OrderItemRes, 0, len(o.Items))

	for _, item := range o.Items {
		items = append(items, toOrderItemRes(item))
	}

	return dto.OrderRes{
		ID:            o.ID,
		Items:         items,
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o dto.OrderReq
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	om := toOrderModel(o)
	om.UserID = claims.ID

	order, err := h.service.CreateOrder(r.Context(), om)
	if err != nil {
		http.Error(w, "error creating order", http.StatusInternalServerError)
		return
	}

	res := toOrderRes(order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	lo, err := h.service.ListOrders(r.Context())
	if err != nil {
		http.Error(w, "error listing orders", http.StatusInternalServerError)
		return
	}

	res := make([]dto.OrderRes, 0, len(lo))
	for _, o := range lo {
		res = append(res, toOrderRes(o))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// ListMyOrders returns the logged-in user's orders
func (h *Handler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	lo, err := h.service.ListOrdersByUser(r.Context(), claims.ID)
	if err != nil {
		http.Error(w, "error listing orders", http.StatusInternalServerError)
		return
	}

	res := make([]dto.OrderRes, 0, len(lo))
	for _, o := range lo {
		res = append(res, toOrderRes(o))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}

	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	o, err := h.service.GetOrder(r.Context(), i)
	if err != nil {
		http.Error(w, "cant get order by this id", http.StatusInternalServerError)
		return
	}

	if !canAccessOrder(claims, o) {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	res := toOrderRes(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}

	claims, ok := claimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	o, err := h.service.GetOrder(r.Context(), i)
	if err != nil {
		http.Error(w, "cant get order by this id", http.StatusInternalServerError)
		return
	}

	if !canAccessOrder(claims, o) {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	err = h.service.DeleteOrder(r.Context(), i)
	if err != nil {
		http.Error(w, "error deleting order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// admins can access any order, users only their own
func canAccessOrder(claims *auth.Claims, o *model.Order) bool {
	return claims.IsAdmin || o.UserID == claims.ID
}
