package handler

import (
	"net/http"
	"strings"

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

// Middleware для проверки авторизации
func (h *Handler) AuthMiddleware(c fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return h.NewErrorResponse(c, http.StatusUnauthorized, "пустой заголовок авторизации")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return h.NewErrorResponse(c, http.StatusUnauthorized, "неверный формат токена")
	}

	userId, err := h.services.Auth.ParseJWT(headerParts[1])
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, "недействительный токен")
	}

	c.Locals("userId", userId)
	return c.Next()
}

// Middleware для проверки прав администратора
func (h *Handler) AdminMiddleware(c fiber.Ctx) error {
	if err := h.AuthMiddleware(c); err != nil {
		return err
	}

	userId := c.Locals("userId").(int)
	isAdmin, err := h.services.Auth.IsAdmin(userId)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, "ошибка проверки прав")
	}

	if !isAdmin {
		return h.NewErrorResponse(c, http.StatusForbidden, "недостаточно прав")
	}

	return c.Next()
}
