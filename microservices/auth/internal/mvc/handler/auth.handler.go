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
		return h.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	id, err := h.services.Auth.CreateUser(user)
	if err != nil {
		return h.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
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
