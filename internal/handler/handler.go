package handler

import "github.com/eunice99x/goMicro/grpc/pb"

// type Handler struct {
// 	service Services
// }

// func NewHandler(service Services) *Handler {
// 	return &Handler{
// 		service: service,
// 	}
// }


type Handler struct {
	client pb.EcommClient
}

func NewHandler(client pb.EcommClient) *Handler {
	return &Handler{
		client: client,
	}
}
