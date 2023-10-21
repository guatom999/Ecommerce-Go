package ordersUseCases

import (
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
)

type IOrderUseCase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
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

func (u *orderUseCase) FindOneOrder(orderId string) (*orders.Order, error) {

	order, err := u.orderRepo.FindOneProduct(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
