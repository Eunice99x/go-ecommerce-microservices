package handler

import (
	"github.com/eunice99x/goMicro/internal/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: *service,
	}
}