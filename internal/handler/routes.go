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

	r.Route("/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.ListUsers)

		// GET   /users/user?email=x@example.com
		// PATCH /users/user?email=x@example.com
		r.Route("/user", func(r chi.Router) {
			r.Get("/", handler.GetUser)
			r.Patch("/", handler.UpdateUser)
		})

		r.Delete("/{id}", handler.DeleteUser)
	})

	// auth
	r.Post("/login", handler.LoginUser)
	r.Post("/refresh", handler.RenewAccessToken)
	r.Delete("/logout/{id}", handler.LogoutUser)

	// sessions
	r.Patch("/sessions/{id}/revoke", handler.RevokeSession)

	return r
}

func Start(addr string, r http.Handler) error {
	return http.ListenAndServe(addr, r)
}
