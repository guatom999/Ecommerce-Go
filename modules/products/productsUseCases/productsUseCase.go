package productsUseCases

import (
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
)

type IProductsUseCase interface {
	FindOneProduct(productId string) (*products.Product, error)
}

type productsUseCase struct {
	productsRepo productsRepositories.IProductRepository
}

func ProductsUseCase(productsRepo productsRepositories.IProductRepository) IProductsUseCase {
	return &productsUseCase{
		productsRepo: productsRepo,
	}
}

func (u *productsUseCase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productsRepo.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}
