package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eunice99x/goMicro/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating m db: %v", err)
	}
	defer mockDB.Close()

	newDB := sqlx.NewDb(mockDB, "sqlmock")
	fn(newDB, m)
}

func TestCreateProduct(t *testing.T) {

	p := &model.Product{
		Name: "test product",
		Image: "test image",
		Category: "test category",
		Description: "test description",
		Rating: 5,
		NumReviews: 10,
		Price: 22,
		CountInStock: 999,
	} 

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name:"sucess",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock)  {

				m.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1,1))
				cp, err := s.CreateProduct(context.Background(), p)
				if err != nil {
					t.Fatalf("error creating product: %v", err)
				}

				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				m.ExpectationsWereMet()
				require.NoError(t, err)	
			},
			
		},
		{
			name: "failed getting last insert ID",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error getting last insert ID")))
				_, err := s.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting product",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnError(fmt.Errorf("error inserting product"))
				_, err := s.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func (t *testing.T)  {
				withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock){
				s := NewPostgresStorer(db)
				tc.test(t, s, m)
			})
		})
	}
}

func TestGetProduct(t *testing.T) {

	p := &model.Product{
		Name: "test product",
		Image: "test image",
		Category: "test category",
		Description: "test description",
		Rating: 5,
		NumReviews: 10,
		Price: 22,
		CountInStock: 999,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name:"sucess",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)

				m.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)

				gp, err := s.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},{
			name: "failed getting product",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting product"))

				_, err := s.GetProduct(context.Background(), 1)
				require.Error(t, err)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				s := NewPostgresStorer(db)
				tc.test(t, s, m)
			})
		})
	}
}

func TestListProducts(t *testing.T) {
	p := &model.Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   100,
		Price:        99.99,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id",
					"name",
					"image",
					"category",
					"description",
					"rating",
					"num_reviews",
					"price",
					"count_in_stock",
					"created_at",
					"updated_at",
				}).AddRow(
					1,
					p.Name,
					p.Image,
					p.Category,
					p.Description,
					p.Rating,
					p.NumReviews,
					p.Price,
					p.CountInStock,
					p.CreatedAt,
					p.UpdatedAt,
				)

				m.ExpectQuery("SELECT * FROM products").
					WillReturnRows(rows)

				products, err := s.ListProducts(context.Background())

				require.NoError(t, err)
				require.Len(t, products, 1)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying products",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT * FROM products").
					WillReturnError(fmt.Errorf("error querying products"))

				_, err := s.ListProducts(context.Background())

				require.Error(t, err)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				s := NewPostgresStorer(db)
				tc.test(t, s, m)
			})
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	p := &model.Product{
		ID:           1,
		Name:         "updated product",
		Image:        "updated.jpg",
		Category:     "updated category",
		Description:  "updated description",
		Rating:       4,
		NumReviews:   50,
		Price:        49.99,
		CountInStock: 20,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec(
					"UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?",
				).WillReturnResult(sqlmock.NewResult(0, 1))

				updatedProduct, err := s.UpdateProduct(
					context.Background(),
					p,
				)

				require.NoError(t, err)
				require.Equal(t, int64(1), updatedProduct.ID)
				require.Equal(t, "updated product", updatedProduct.Name)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating product",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec(
					"UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?",
				).WillReturnError(fmt.Errorf("error updating product"))

				_, err := s.UpdateProduct(
					context.Background(),
					p,
				)

				require.Error(t, err)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				s := NewPostgresStorer(db)
				tc.test(t, s, m)
			})
		})
	}
}


func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM products WHERE id=?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := s.DeleteProduct(context.Background(), 1)

				require.NoError(t, err)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T, s *PostgresStorer, m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM products WHERE id=?").
					WithArgs(1).
					WillReturnError(fmt.Errorf("error deleting product"))

				err := s.DeleteProduct(context.Background(), 1)

				require.Error(t, err)

				err = m.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, m sqlmock.Sqlmock) {
				s := NewPostgresStorer(db)
				tc.test(t, s, m)
			})
		})
	}
}