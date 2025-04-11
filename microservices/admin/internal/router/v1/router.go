package v1

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/Linkr/admin/internal/handler"
)

func NewRouter(h *handler.Handler) *fiber.App {
	app := fiber.New()

	// Группа для административных функций
	admin := app.Group("/api/admin/v1")
	{
		admin.Get("/users", h.GetUsers)
		admin.Get("/users/:id", h.GetUserById)
		admin.Put("/users/:id/block", h.BlockUser)
		admin.Put("/users/:id/role", h.ChangeUserRole)
	}

	return app
}
