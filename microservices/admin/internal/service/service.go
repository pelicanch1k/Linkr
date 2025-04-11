package service

import (
	"github.com/pelicanch1k/Linkr/admin/internal/dto"
	"github.com/pelicanch1k/Linkr/admin/internal/repository"
)

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
	Admin
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		NewAdminService(repos.Admin),
	}
}
