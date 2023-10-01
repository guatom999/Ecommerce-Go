package middlewaresHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresUsecases"
	"github.com/guatom999/Ecommerce-Go/pkg/authen"
	"github.com/guatom999/Ecommerce-Go/pkg/utils"
)

type middlewareHandlerErrorCode string

const (
	routerCheckErr middlewareHandlerErrorCode = "middleware-001"
	jwtAuthErr     middlewareHandlerErrorCode = "middleware-002"
	paramsCheckErr middlewareHandlerErrorCode = "middleware-003"
	authorizeErr   middlewareHandlerErrorCode = "middleware-004"
	apiKeyErr      middlewareHandlerErrorCode = "middleware-005"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoldId ...int) fiber.Handler
	ApiKeyCheck() fiber.Handler
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

func (h *middlewareHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		result, err := authen.ParseToken(h.cfg.Jwt(), token)

		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims

		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewareHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"params doesn't match",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewareHandler) Authorize(expectRoldId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}

		roles, err := h.middlewaresUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0

		for _, v := range expectRoldId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if expectedValueBinary[i]&userValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()

		// // fmt.Printf("expectedValueBinary :%v", expectedValueBinary)
		// // fmt.Printf("userValueBinary :%v", userValueBinary)

		// for i := range userValueBinary {
		// 	fmt.Printf("compare value is :%v", expectedValueBinary[i]&userValueBinary[i])
		// 	if expectedValueBinary[i]&userValueBinary[i] != 1 {
		// 		return entities.NewResponse(c).Error(
		// 			fiber.ErrUnauthorized.Code,
		// 			string(authorizeErr),
		// 			"no permission to access",
		// 		).Res()
		// 	}
		// }

		// return c.Next()

	}
}

func (h *middlewareHandler) ApiKeyCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {

		key := c.Get("X-Api-Key")
		if _, err := authen.ParseApikey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(apiKeyErr),
				"apikey is invalid",
			).Res()
		}

		return c.Next()
	}
}
