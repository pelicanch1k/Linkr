package dto

// Регистрация
type RegisterUser struct {
	Username   string `json:"username" db:"username" binding:"required"`
	Password   string `json:"password" db:"password_hash" binding:"required"`
	Email      string `json:"email" db:"email" binding:"required,email"`
	First_name string `json:"first_name" db:"first_name" binding:"required"`
	Last_name  string `json:"last_name" db:"last_name" binding:"required"`
}

// Авторизация
type AuthorizationUser struct {
	Username string `json:"username" db:"username" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
}

// Работа с токенами
type RefreshToken struct {
	Token string `json:"token" binding:"required"`
}

// Восстановление пароля
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
