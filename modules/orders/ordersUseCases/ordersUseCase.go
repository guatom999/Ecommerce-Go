package ordersUseCases

import (
	"fmt"
	"math"

	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
)

type IOrderUseCase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
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

func (u *orderUseCase) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {

	orders, count := u.orderRepo.FindOrder(req)

	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *orderUseCase) InsertOrder(req *orders.Order) (*orders.Order, error) {

	for i := range req.Product {
		if req.Product[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}

		product, err := u.productRepo.FindOneProduct(req.Product[i].Id)

		if err != nil {
			return nil, err
		}

		req.TotalPaid += req.Product[i].Product.Price * float64(req.Product[i].Quantity)
		req.Product[i].Product = product

	}

	orderId, err := u.orderRepo.InsertOrder(req)

	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.FindOneProduct(orderId)

	if err != nil {
		return nil, err
	}

	return order, nil
}
