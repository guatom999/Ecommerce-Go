package ordersUseCases

import (
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
)

type IOrderUseCase interface {
}

type orderUseCase struct {
	orderRepo   ordersRepositories.IOrderRepository
	productRepo productsRepositories.IProductRepository
}

func OrderUseCase(orderRepo ordersRepositories.IOrderRepository, productRepo productsRepositories.IProductRepository) IOrderUseCase {
	return &orderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}
