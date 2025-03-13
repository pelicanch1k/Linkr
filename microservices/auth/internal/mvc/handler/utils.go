package handler

import "github.com/gofiber/fiber/v3"

func (h *Handler) NewErrorResponse(c fiber.Ctx, statusCode int, message string) error {
	h.logger.Error(message)
	return c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}