package handler

import (
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(handler *Handler, tokenGen *auth.JWTConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/products", func(r chi.Router) {
		r.With(GetAdminMiddlewareFunc(tokenGen)).Post("/", handler.CreateProduct)
		r.Get("/", handler.ListProducts)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetProduct)
			r.Group(func(r chi.Router) {
				r.Use(GetAdminMiddlewareFunc(tokenGen))
				r.Patch("/", handler.UpdateProduct)
				r.Delete("/", handler.DeleteProduct)
			})
		})
	})

	r.Route("/orders", func(r chi.Router) {
		r.Use(GetAuthMiddlewareFunc(tokenGen))

		r.Post("/", handler.CreateOrder)
		r.With(GetAdminMiddlewareFunc(tokenGen)).Get("/", handler.ListOrders)
		r.Get("/me", handler.ListMyOrders)

		// owner or admin, checked in the handler
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetOrder)
			// r.Patch("/", handler.UpdateOrder)
			r.Delete("/", handler.DeleteOrder)
		})
	})

	r.Route("/users", func(r chi.Router) {
		// public
		r.Post("/", handler.CreateUser)

		// admin only
		r.Group(func(r chi.Router) {
			r.Use(GetAdminMiddlewareFunc(tokenGen))

			r.Get("/", handler.ListUsers)
			r.Delete("/{id}", handler.DeleteUser)
		})

		// authenticated user
		r.Group(func(r chi.Router) {
			r.Use(GetAuthMiddlewareFunc(tokenGen))

			r.Route("/user", func(r chi.Router) {
				// GET   /users/user              (self; admins may pass ?email=x@example.com)
				// PATCH /users/user              (self)
				r.Get("/", handler.GetUser)
				r.Patch("/", handler.UpdateUser)
			})
		})
	})

	// auth - refresh is public since the access token is usually expired by then
	r.Post("/login", handler.LoginUser)
	r.Post("/refresh", handler.RenewAccessToken)

	// authenticated token/session operations
	r.Group(func(r chi.Router) {
		r.Use(GetAuthMiddlewareFunc(tokenGen))

		// logs out the current session (from the token's sid)
		r.Delete("/logout", handler.LogoutUser)
		// revoke any of your own sessions, e.g. "log out other devices"
		r.Patch("/sessions/{id}/revoke", handler.RevokeSession)
	})

	return r
}
