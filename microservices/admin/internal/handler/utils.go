package handler

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewErrorResponse(c fiber.Ctx, statusCode int, message string) error {
	h.logger.Error(message)
	return c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}

func (h *Handler) NewSuccessResponse(c fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"data": data,
	})
}

// Получение ID пользователя из контекста
func (h *Handler) GetUserIdFromContext(c fiber.Ctx) (int, error) {
	userId, ok := c.Locals("userId").(int)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "пользователь не авторизован")
	}
	return userId, nil
}

// Проверка валидности email
func (h *Handler) ValidateEmail(email string) bool {
	// Простая проверка на наличие @ и домена
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Проверка валидности пароля
func (h *Handler) ValidatePassword(password string) bool {
	return len(password) >= 6
}
