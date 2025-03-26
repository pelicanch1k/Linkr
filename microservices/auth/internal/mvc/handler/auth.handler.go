package handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
)

func (h *Handler) SignUp(c fiber.Ctx) error {
	var user dto.RegisterUser

	if err := c.Bind().Body(&user); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, "Ошибка в данных формы: "+err.Error())
	}

	id, err := h.services.Auth.CreateUser(user)
	if err != nil {
		statusCode := http.StatusInternalServerError
		// Определяем различные ошибки валидации для более точного статус-кода
		if strings.Contains(err.Error(), "email") ||
			strings.Contains(err.Error(), "имя") ||
			strings.Contains(err.Error(), "фамилия") ||
			strings.Contains(err.Error(), "уже существует") {
			statusCode = http.StatusBadRequest
		}
		return h.NewErrorResponse(c, statusCode, err.Error())
	}

	return c.Status(http.StatusOK).JSON(map[string]int{
		"id": id,
	})
}

func (h *Handler) SignIn(c fiber.Ctx) error {
	var user dto.AuthorizationUser

	if err := c.Bind().Body(&user); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.Auth.GenerateJWT(user)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.logger.Info("token: ", token)

	return c.Status(http.StatusOK).JSON(map[string]string{
		"token": token,
	})
}

func (h *Handler) UserIdentity(c fiber.Ctx) error {
	header := c.Get("Authorization")

	if header == "" {
		return h.NewErrorResponse(c, http.StatusBadRequest, "empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || len(headerParts[1]) == 0 {
		return h.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
	}

	userId, err := h.services.Auth.ParseJWT(headerParts[1])
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	h.logger.Info(userId)

	return c.Status(http.StatusOK).JSON(map[string]int{
		"user_id": userId,
	})
}

func (h *Handler) SignOut(c fiber.Ctx) error {
	userId, err := h.GetUserIdFromContext(c)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	if err := h.services.Auth.Logout(userId); err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "успешный выход",
	})
}

func (h *Handler) RefreshToken(c fiber.Ctx) error {
	var refreshToken dto.RefreshToken
	if err := c.Bind().Body(&refreshToken); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	newToken, err := h.services.Auth.RefreshToken(refreshToken.Token)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": newToken,
	})
}

func (h *Handler) ForgotPassword(c fiber.Ctx) error {
	var request dto.ForgotPasswordRequest
	if err := c.Bind().Body(&request); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Auth.SendResetPasswordEmail(request.Email); err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "инструкции по сбросу пароля отправлены на email",
	})
}

func (h *Handler) ResetPassword(c fiber.Ctx) error {
	var request dto.ResetPasswordRequest
	if err := c.Bind().Body(&request); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Auth.ResetPassword(request.Token, request.NewPassword); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "пароль успешно изменен",
	})
}

func (h *Handler) ChangePassword(c fiber.Ctx) error {
	userId, err := h.GetUserIdFromContext(c)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	var request dto.ChangePasswordRequest
	if err := c.Bind().Body(&request); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Auth.ChangePassword(userId, request.OldPassword, request.NewPassword); err != nil {
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "пароль успешно изменен",
	})
}
