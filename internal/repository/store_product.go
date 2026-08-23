package repository

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
)

func (ps *PostgresStorer) CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	query := `
		INSERT INTO products (
			name,
			image,
			category,
			description,
			rating,
			num_reviews,
			price,
			count_in_stock
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := ps.db.GetContext(
		ctx,
		&p.ID,
		query,
		p.Name,
		p.Image,
		p.Category,
		p.Description,
		p.Rating,
		p.NumReviews,
		p.Price,
		p.CountInStock,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", err)
	}

	return p, nil
}

func (ps *PostgresStorer) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	var p model.Product

	err := ps.db.GetContext(
		ctx,
		&p,
		"SELECT * FROM products WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	return &p, nil
}

func (ps *PostgresStorer) ListProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product

	err := ps.db.SelectContext(
		ctx,
		&products,
		"SELECT * FROM products",
	)
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return products, nil
}

func (ps *PostgresStorer) UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	query := `
		UPDATE products
		SET
			name=$1,
			image=$2,
			category=$3,
			description=$4,
			rating=$5,
			num_reviews=$6,
			price=$7,
			count_in_stock=$8,
			updated_at=$9
		WHERE id=$10
	`

	_, err := ps.db.ExecContext(
		ctx,
		query,
		p.Name,
		p.Image,
		p.Category,
		p.Description,
		p.Rating,
		p.NumReviews,
		p.Price,
		p.CountInStock,
		p.UpdatedAt,
		p.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return p, nil
}

func (ps *PostgresStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := ps.db.ExecContext(
		ctx,
		"DELETE FROM products WHERE id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}