package api

import (
	"context"
	"fmt"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/jmoiron/sqlx"
)

type PostgresStorer struct {
	db *sqlx.DB
}

func NewPostgresStorer(db *sqlx.DB) *PostgresStorer {
	return &PostgresStorer{
		db:db,
	}
}

func (ps *PostgresStorer) CreateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	res, err := ps.db.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)
	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last id: %v", id)
	}
	p.ID = id

	return p, nil
}

func (ps *PostgresStorer) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	var p model.Product
	err := ps.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	return &p, nil
}

func (ps *PostgresStorer) ListProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	err := ps.db.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return products, nil
}

func (ps *PostgresStorer) UpdateProduct(ctx context.Context, p *model.Product) (*model.Product, error) {
	_, err := ps.db.NamedExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock, updated_at=:updated_at WHERE id=:id", p)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return p, nil
}

func (ps *PostgresStorer) DeleteProduct(ctx context.Context, id int64) (error) {
	_, err := ps.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}