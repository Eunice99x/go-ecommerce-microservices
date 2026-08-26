package service

import (
	"context"

	"github.com/eunice99x/goMicro/internal/model"
)


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