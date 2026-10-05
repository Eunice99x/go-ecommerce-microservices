package service

import (
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/internal/repository"
)


type Service struct {
	storer *repository.PostgresStorer
	pb.UnimplementedEcommServer
}

func NewStorer(storer *repository.PostgresStorer) *Service {
	return &Service{
		storer: storer,
	}
}