package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/Linkr/admin/internal/dto"
)

func (h *Handler) GetUsers(c fiber.Ctx) error {
	h.logger.Info(c.OriginalURL())

	users, err := h.services.GetUsers()
	if err != nil {
		h.logger.Error("Ошибка получения пользователей: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.logger.Info("Пользователи получены успешно")
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": users,
	})
}

func (h *Handler) GetUserById(c fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, "неверный id пользователя")
	}

	user, err := h.services.GetUserById(userId)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, user)
}

func (h *Handler) BlockUser(c fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, "неверный id пользователя")
	}

	var request dto.BlockUserRequest
	if err := c.Bind().Body(&request); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.BlockUser(userId, request.Blocked); err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, fiber.Map{"message": "статус блокировки обновлен"})
}

func (h *Handler) ChangeUserRole(c fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, "неверный id пользователя")
	}

	var request dto.ChangeRoleRequest
	if err := c.Bind().Body(&request); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.ChangeUserRole(userId, request.Role); err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, fiber.Map{"message": "роль пользователя изменена"})
}

