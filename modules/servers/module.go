package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
}

// Constructor
func NewModule(router fiber.Router, server *server) IModuleFactory {
	return &moduleFactory{
		router: router,
		server: server,
	}
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}
