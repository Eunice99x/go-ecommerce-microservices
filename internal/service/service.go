package service

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/eunice99x/goMicro/internal/repository"
)

type Service struct {
	storer *repository.PostgresStorer
}

func NewService(storer *repository.PostgresStorer) *Service {
	return &Service{storer: storer}
}

func (s *Service) CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error){
	return s.storer.CreateProduct(ctx, p)
}

func (s *Service) GetProduct(ctx context.Context, id int64) (*model.Product, error){
	return s.storer.GetProduct(ctx, id)
}

func (s *Service) ListProducts(ctx context.Context) ([]*model.Product, error){
	return s.storer.ListProducts(ctx)
}

func (s *Service) UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	return s.storer.UpdateProduct(ctx, p)
}

func (s *Service) DeleteProduct(ctx context.Context, id int64) error {
	return s.storer.DeleteProduct(ctx, id)
}