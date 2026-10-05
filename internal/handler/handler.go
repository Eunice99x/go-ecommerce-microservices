package handler

type Handler struct {
	service Services
}

func NewHandler(service Services) *Handler {
	return &Handler{
		service: service,
	}
}
