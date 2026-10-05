package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/eunice99x/goMicro/internal/model"
)

// writeServiceError maps domain errors to http statuses; unknown errors become a plain 500
func writeServiceError(w http.ResponseWriter, err error, resource string) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		http.Error(w, resource+" not found", http.StatusNotFound)
	case errors.Is(err, model.ErrAlreadyExists):
		http.Error(w, resource+" already exists", http.StatusConflict)
	case errors.Is(err, model.ErrInvalidArgument):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, model.ErrInvalidStatusTransition):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, model.ErrInvalidCredentials):
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
	default:
		log.Printf("handler: %s: %v", resource, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
