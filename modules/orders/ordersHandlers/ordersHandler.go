package ordersHandlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersUseCases"
)

type orderErrCode string

const (
	findeOneOrderErr orderErrCode = "order-001"
	findOrderErr     orderErrCode = "order-002"
	insertOrderErr   orderErrCode = "order-003"
	updateOrderErr   orderErrCode = "order-004"
)

type IOrderHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
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

func (h *orderHandler) InsertOrder(c *fiber.Ctx) error {

	userId := c.Locals("userId").(string)

	req := &orders.Order{
		Product: make([]*orders.ProductsOrder, 0),
	}

	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	if len(req.Product) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			"product are empthy",
		).Res()
	}

	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.orderUseCase.InsertOrder(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		order,
	).Res()
}

func (h *orderHandler) UpdateOrder(c *fiber.Ctx) error {

	orderId := strings.Trim(c.Params("order_id"), " ")
	req := new(orders.Order)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}
	req.Id = orderId

	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled",
	}
	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else if strings.ToLower(req.Status) == statusMap["canceled"] {
		req.Status = statusMap["canceled"]
	} else {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateOrderErr),
			"customer cannot do that",
		).Res()
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}
		if req.TransferSlip.CreatedAt == "" {
			localTime, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(updateOrderErr),
					err.Error(),
				).Res()
			}

			now := time.Now().In(localTime)

			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")
		}

	}

	order, err := h.orderUseCase.UpdateOrder(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		order,
	).Res()
}
