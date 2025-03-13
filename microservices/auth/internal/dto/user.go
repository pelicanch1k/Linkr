package dto

type RegisterUser struct {
	Username string `json:"username" db:"username" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
	Email string `json:"email" db:"email" binding:"required"`
	First_name string `json:"first_name" db:"first_name" binding:"required"`
	Last_name string `json:"last_name" db:"last_name" binding:"required"`
}

type AuthorizationUser struct {
	Username string `json:"username" db:"username" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
}
