package ordersHandlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersUseCases"
)

type orderErrCode string

const (
	findeOneOrderErr orderErrCode = "order-001"
	findOrderErr     orderErrCode = "order-002"
)

type IOrderHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
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

func (h *orderHandler) FindOrder(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	// Paginate
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Sort
	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date	YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"start date is invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}
	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"end date is invalid",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		h.orderUseCase.FindOrder(req),
	).Res()
}
