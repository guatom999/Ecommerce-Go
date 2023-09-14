package usersHandlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/users"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersUsecases"
)

type userHandlersErrCode string

const (
	signUpCustomerErr userHandlersErrCode = "users-001"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

func UsersHandler(cfg config.IConfig, userUsecase usersUsecases.IUserUsecase) IUsersHandler {
	return &usersHandler{cfg, userUsecase}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {

	// body parser
	req := new(users.UserRegisterReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// validation email

	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern doesn't match",
		).Res()
	}

	// Insert User
	result, err := h.userUsecase.InsertCustomer(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			fmt.Sprintf("some data in used :%v", err),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
