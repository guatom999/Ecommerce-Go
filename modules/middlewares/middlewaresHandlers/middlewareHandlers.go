package middlewaresHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresUsecases"
)

type middlewareHandlerErrorCode string

const (
	routerCheckErr middlewareHandlerErrorCode = "middleware-001"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
}

type middlewareHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewareHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewareHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewareHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewareHandler) Logger() fiber.Handler {
	// return func(c *fiber.Ctx) error {}
	return logger.New(logger.Config{
		Format:     "${time} [${ip}]  ${status} - ${method} ${path}\n",
		TimeFormat: "01/02/2006",
		TimeZone:   "Bangkok/Asia",
	})
}
