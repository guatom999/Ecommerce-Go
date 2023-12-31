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
	AddProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
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

func (u *productsUseCase) AddProduct(req *products.Product) (*products.Product, error) {

	product, err := u.productsRepo.InsertProduct(req)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productsUseCase) UpdateProduct(req *products.Product) (*products.Product, error) {

	product, err := u.productsRepo.UpdateProduct(req)

	if err != nil {
		return nil, err
	}

	return product, nil

}

func (u *productsUseCase) DeleteProduct(productId string) error {

	if err := u.productsRepo.DeleteProduct(productId); err != nil {
		return err
	}

	return nil
}
