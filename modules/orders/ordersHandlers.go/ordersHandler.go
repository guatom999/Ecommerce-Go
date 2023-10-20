package ordersHandlers

import (
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersUseCases.go"
)

type IOrderHandler interface {
}

type orderHandler struct {
	config       config.IConfig
	orderUseCase ordersUseCases.IOrderUseCase
}

func OrderHandler(config config.IConfig, orderUseCase ordersUseCases.IOrderUseCase) IOrderHandler {
	return &orderHandler{
		config:       config,
		orderUseCase: orderUseCase,
	}
}
