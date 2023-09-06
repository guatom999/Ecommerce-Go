package monitorHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/monitor"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (m *monitorHandler) HealthCheck(c *fiber.Ctx) error {

	res := &monitor.Monitor{
		Name:    m.cfg.App().Name(),
		Version: m.cfg.App().Version(),
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
