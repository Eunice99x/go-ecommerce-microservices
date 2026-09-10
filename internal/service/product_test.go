package service

import (
	"fmt"
	"testing"

	"github.com/eunice99x/goMicro/internal/model"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:           1,
					Name:         "Iphone",
					Image:        "https://example.com",
					Category:     "Electronics",
					Price:        999,
					CountInStock: 1234,
				}

				s := &Service{
					storer: &fakeStorer{product: p},
				}

				got, err := s.CreateProduct(t.Context(), p)

				require.NoError(t, err)
				require.Equal(t, p, got)
			},
		},
		{
			name: "failed creating product",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error creating product")},
				}

				_, err := s.CreateProduct(t.Context(), &model.Product{})

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestGetProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:    1,
					Name:  "Iphone",
					Price: 999,
				}

				s := &Service{
					storer: &fakeStorer{product: p},
				}

				got, err := s.GetProduct(t.Context(), 1)

				require.NoError(t, err)
				require.Equal(t, int64(1), got.ID)
				require.Equal(t, "Iphone", got.Name)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error getting product")},
				}

				_, err := s.GetProduct(t.Context(), 1)

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestListProducts(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				products := []*model.Product{
					{ID: 1, Name: "Iphone", Price: 999},
					{ID: 2, Name: "Macbook", Price: 1500},
				}

				s := &Service{
					storer: &fakeStorer{products: products},
				}

				got, err := s.ListProducts(t.Context())

				require.NoError(t, err)
				require.Len(t, got, 2)
				require.Equal(t, "Iphone", got[0].Name)
				require.Equal(t, "Macbook", got[1].Name)
			},
		},
		{
			name: "failed listing products",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error listing products")},
				}

				_, err := s.ListProducts(t.Context())

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				p := &model.Product{
					ID:    1,
					Name:  "Updated Iphone",
					Price: 1200,
				}

				s := &Service{
					storer: &fakeStorer{product: p},
				}

				got, err := s.UpdateProduct(t.Context(), p)

				require.NoError(t, err)
				require.Equal(t, "Updated Iphone", got.Name)
				require.Equal(t, float64(1200), got.Price)
			},
		},
		{
			name: "failed updating product",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error updating product")},
				}

				_, err := s.UpdateProduct(t.Context(), &model.Product{})

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{},
				}

				err := s.DeleteProduct(t.Context(), 1)

				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T) {
				s := &Service{
					storer: &fakeStorer{err: fmt.Errorf("error deleting product")},
				}

				err := s.DeleteProduct(t.Context(), 1)

				require.Error(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
