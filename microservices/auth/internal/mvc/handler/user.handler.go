package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
)

func (h *Handler) GetProfile(c fiber.Ctx) error {
	userId, err := h.GetUserIdFromContext(c)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	profile, err := h.services.User.GetProfile(userId)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	userId, err := h.GetUserIdFromContext(c)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var request dto.UpdateProfileRequest
	if err := c.Bind().Body(&request); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.User.UpdateProfile(userId, request); err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, fiber.Map{"message": "профиль обновлен"})
}

func (h *Handler) DeleteProfile(c fiber.Ctx) error {
	userId, err := h.GetUserIdFromContext(c)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	if err := h.services.User.DeleteProfile(userId); err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, fiber.Map{"message": "профиль удален"})
}

func (h *Handler) UpdateAvatar(c fiber.Ctx) error {
	userId, err := h.GetUserIdFromContext(c)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, "ошибка загрузки файла")
	}

	avatarURL, err := h.services.User.UpdateAvatar(userId, file)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.NewSuccessResponse(c, http.StatusOK, fiber.Map{"avatar_url": avatarURL})
}
