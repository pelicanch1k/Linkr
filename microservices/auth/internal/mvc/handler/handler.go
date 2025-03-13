package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/ProductGatewayAPI/pkg/logging"

	"github.com/pelicanch1k/Linkr/auth/internal/mvc/service"
)


type Handler struct {
	services *service.Service
	logger   *logging.Logger
}

func NewHandler(services *service.Service, logger *logging.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) Logging(c fiber.Ctx) error {
	h.logger.Info(c.BaseURL())
	return c.Next()
}