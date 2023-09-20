package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresUsecases"
	"github.com/guatom999/Ecommerce-Go/modules/monitor/monitorHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.IMiddlewareHandler
}

// Constructor
func NewModule(router fiber.Router, server *server, mid middlewaresHandlers.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router: router,
		server: server,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewareHandler {
	repository := middlewaresRepositories.MiddlewareRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	handler := middlewaresHandlers.MiddlewareHandler(s.cfg, usecase)

	return handler
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefeshPassport)

}
