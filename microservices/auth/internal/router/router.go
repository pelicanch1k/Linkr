package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/handler"
)

func NewRouter(h *handler.Handler) *fiber.App {
	app := fiber.New()

	auth := app.Group("/api/auth/v1", h.Logging)
	{
		auth.Post("/sign-up", h.SignUp)
		auth.Post("/sign-in", h.SignIn)
		auth.Post("/check-user", h.UserIdentity)
	}

	return app
}