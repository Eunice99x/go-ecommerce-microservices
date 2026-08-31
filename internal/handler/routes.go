package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/products", func(r chi.Router) {
		r.Post("/", handler.CreateProduct)
		r.Get("/", handler.ListProducts)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetProduct)
			r.Patch("/", handler.UpdateProduct)
			r.Delete("/", handler.DeleteProduct)
		})
	})

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", handler.CreateOrder)
		r.Get("/", handler.ListOrders)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetOrder)
			// r.Patch("/", handler.UpdateOrder)
			r.Delete("/", handler.DeleteOrder)
		})
	})

	return r
}

func Start(addr string, r http.Handler) error {
	return http.ListenAndServe(addr, r)
}
