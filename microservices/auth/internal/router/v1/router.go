package v1

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/handler"
)

func NewRouter(h *handler.Handler) *fiber.App {
	app := fiber.New()

	// Группа для аутентификации
	auth := app.Group("/api/auth/v1", h.Logging)
	{
		auth.Post("/sign-up", h.SignUp)
		auth.Post("/sign-in", h.SignIn)
		auth.Post("/sign-out", h.SignOut, h.AuthMiddleware)
		auth.Post("/refresh-token", h.RefreshToken)
		auth.Post("/check-user", h.UserIdentity)
		auth.Post("/forgot-password", h.ForgotPassword)
		auth.Post("/reset-password", h.ResetPassword)
		auth.Put("/change-password", h.ChangePassword)
	}

	// Группа для управления профилем пользователя
	user := app.Group("/api/users/v1", h.AuthMiddleware)
	{
		user.Get("/me", h.GetProfile)
		user.Put("/me", h.UpdateProfile)
		user.Delete("/me", h.DeleteProfile)
		user.Put("/me/avatar", h.UpdateAvatar)
	}

	// Группа для административных функций
	// admin := app.Group("/api/admin/v1", h.AdminMiddleware)
	// {
	// 	admin.Get("/users", h.GetUsers)
	// 	admin.Get("/users/:id", h.GetUserById)
	// 	admin.Put("/users/:id/block", h.BlockUser)
	// 	admin.Put("/users/:id/role", h.ChangeUserRole)
	// }

	return app
}
