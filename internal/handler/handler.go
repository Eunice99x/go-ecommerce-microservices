package handler

import (
	"context"

	"github.com/eunice99x/goMicro/internal/service"
)

type Handler struct {
	ctx context.Context
	service service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		ctx: context.Background(),
		service: *service,
	}
}