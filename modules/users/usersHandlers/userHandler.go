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
	signUpCustomerErr  userHandlersErrCode = "users-001"
	signInCustomerErr  userHandlersErrCode = "users-002"
	refreshPassportErr userHandlersErrCode = "users-003"
	signOutErr         userHandlersErrCode = "users-004"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefeshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
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

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInCustomerErr),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecase.GetPassport(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInCustomerErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefeshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)

	if err := c.BodyParser(req); err != nil {
		fmt.Sprintln("Here 1")
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecase.RefreshPassport(req)

	if err != nil {
		fmt.Sprintln("Here 1")
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {

	req := new(users.UserRemoveCredential)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.userUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		nil,
	).Res()
}
