package ordersHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersUseCases.go"
)

type orderErrCode string

const (
	findeOneOrderErr orderErrCode = "order-001"
)

type IOrderHandler interface {
	FindOneOrder(c *fiber.Ctx) error
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

func (h *orderHandler) FindOneOrder(c *fiber.Ctx) error {

	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.orderUseCase.FindOneOrder(orderId)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findeOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		order,
	).Res()
}
