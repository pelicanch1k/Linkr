package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	// "github.com/pelicanch1k/Linkr/auth/internal/model"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository/postgres"
)

type Auth interface {
	CreateUser(user dto.RegisterUser) (int, error)
	GetUserId(user dto.AuthorizationUser) (int, error)
	CheckPassword(userId int, hashedPassword string) (bool, error)
	UpdatePassword(userId int, hashedPassword string) error
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	StoreToken(userId int, token string) error
	CheckToken(userId int, token string) (bool, error)
	InvalidateToken(userId int, token string) error
	InvalidateTokens(userId int) error
	StoreResetToken(email string, token string) error
	ValidateResetToken(token string) (int, error)
	IsAdmin(userId int) (bool, error)
}

type User interface {
	GetUserProfile(userId int) (dto.UserProfile, error)
	UpdateUserProfile(userId int, profile dto.UpdateProfileRequest) error
	DeleteUser(userId int) error
	UpdateAvatar(userId int, avatarURL string) error
	UpdateEmail(userId int, email string) error
	CheckUsernameExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	GetUserNotifications(userId int) ([]dto.Notification, error)
}

type Repository struct {
	Auth
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:  postgres.NewAuthRepository(db),
		User:  postgres.NewUserRepository(db),
	}
}