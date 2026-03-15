package productsRepositories

import (
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/stretchr/testify/mock"
)

type (
	MockProductRepository struct {
		mock.Mock
	}
)

func (m *MockProductRepository) FindOneProduct(productId string) (*products.Product, error) {
	args := m.Called(productId)
	return args.Get(0).(*products.Product), args.Error(1)
}
func (m *MockProductRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	args := m.Called(req)
	return args.Get(0).([]*products.Product), args.Int(1)
}
func (m *MockProductRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	args := m.Called(req)
	return args.Get(0).(*products.Product), args.Error(1)
}
func (m *MockProductRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	args := m.Called(req)
	return args.Get(0).(*products.Product), args.Error(1)
}
func (m *MockProductRepository) DeleteProduct(productId string) error {
	args := m.Called(productId)
	return args.Error(0)
}
