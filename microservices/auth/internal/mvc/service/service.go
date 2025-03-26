package service

import (
	"mime/multipart"

	"github.com/pelicanch1k/Linkr/auth/internal/config"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository"
)

// Auth интерфейс для аутентификации
type Auth interface {
	CreateUser(user dto.RegisterUser) (int, error)
	GenerateJWT(user dto.AuthorizationUser) (string, error)
	ParseJWT(tokenString string) (int, error)
	Logout(userId int) error
	RefreshToken(refreshToken string) (string, error)
	SendResetPasswordEmail(email string) error
	ResetPassword(token, newPassword string) error
	ChangePassword(userId int, oldPassword, newPassword string) error
	IsAdmin(userId int) (bool, error)
	ValidateToken(token string) (bool, error)
	RevokeAllTokens(userId int) error
}

// User интерфейс для работы с профилем
type User interface {
	GetProfile(userId int) (dto.UserProfile, error)
	UpdateProfile(userId int, profile dto.UpdateProfileRequest) error
	DeleteProfile(userId int) error
	UpdateAvatar(userId int, file *multipart.FileHeader) (string, error)
	UpdateEmail(userId int, email string) error
	GetNotifications(userId int) ([]dto.Notification, error)
}

// Admin интерфейс для админских функций
type Admin interface {
	GetUsers() ([]dto.UserProfile, error)
	GetUserById(userId int) (dto.UserProfile, error)
	BlockUser(userId int, blocked bool) error
	ChangeUserRole(userId int, role string) error
	GetUserStats(userId int) (dto.UserStats, error)
	GetSystemStats() (dto.SystemStats, error)
	DeleteUser(userId int) error
}

type Service struct {
	Auth
	User
	Admin
}

func NewService(repos *repository.Repository, authConfig *config.AuthConfig) *Service {
	return &Service{
		Auth:  NewAuthService(repos.Auth, authConfig),
		User:  NewUserService(repos.User),
		Admin: NewAdminService(repos.Admin),
	}
}
