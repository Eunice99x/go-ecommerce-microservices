package service

import (
	"github.com/eunice99x/goMicro/internal/repository"
)

type Service struct {
	storer *repository.PostgresStorer
}

func NewService(storer *repository.PostgresStorer) *Service {
	return &Service{storer: storer}
}
