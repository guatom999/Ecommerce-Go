package productsUseCases

import (
	"math"

	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
)

type IProductsUseCase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) *entities.PaginateRes
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

func (u *productsUseCase) FindProduct(req *products.ProductFilter) *entities.PaginateRes {

	products, count := u.productsRepo.FindProduct(req)

	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
		TotalItem: count,
	}

}
