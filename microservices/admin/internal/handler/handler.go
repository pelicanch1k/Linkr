package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/ProductGatewayAPI/pkg/logging"

	"github.com/pelicanch1k/Linkr/admin/internal/service"
)

type Handler struct {
	services *service.AdminService
	logger   *logging.Logger
}

func NewHandler(services *service.AdminService, logger *logging.Logger) *Handler {
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

	// URL для POST-запроса
	url := "http://127.0.0.1:80/api/auth/v1/check-user"

	// Создание запроса
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		h.logger.Info(fmt.Sprintf("Ошибка при создании запроса: %s", err.Error()))
	}

	// Добавление заголовка Authorization
	req.Header.Set("Authorization", header)
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Info(fmt.Sprintf("Ошибка при отправке запроса: %s", err.Error()))
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.logger.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	// Проверка на успешный ответ
	if resp.StatusCode == http.StatusOK {
		// Парсинг JSON ответа
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			h.logger.Fatalf("Ошибка при парсинге JSON: %v", err)
		}

		userId, ok := response["user_id"]
		if !ok {
			return h.NewErrorResponse(c, http.StatusInternalServerError, "неверный ответ от сервера")
		}

		// Вывод user_id
		h.logger.Info("userId: ", userId)

		c.Locals("userId", userId)
		return c.Next()
	} else {
		h.logger.Fatalf("Ошибка: %s\n", resp.Status)
	}

	return h.NewErrorResponse(c, http.StatusInternalServerError, "неверный ответ от сервера")
}

// Middleware для проверки прав администратора
func (h *Handler) AdminMiddleware(c fiber.Ctx) error {
	if err := h.AuthMiddleware(c); err != nil {
		return err
	}

	userIdVal := c.Locals("userId")
	userIdFloat, ok := userIdVal.(float64)
	if !ok {
		return h.NewErrorResponse(c, http.StatusInternalServerError, "userId имеет неверный тип")
	}

	userId := int(userIdFloat)

	isAdmin, err := h.services.IsAdmin(userId)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, "ошибка проверки прав")
	}

	if !isAdmin {
		return h.NewErrorResponse(c, http.StatusForbidden, "недостаточно прав")
	}

	h.logger.Info("юзер ", userId, " - админ")
	return c.Next()
}
