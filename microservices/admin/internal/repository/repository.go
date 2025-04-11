package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/admin/internal/dto"
	"github.com/pelicanch1k/Linkr/admin/internal/model"
	"github.com/pelicanch1k/Linkr/admin/internal/repository/postgres"
)

type Admin interface {
	GetAllUsers() ([]model.User, error)
	GetUserById(userId int) (model.User, error)
	UpdateUserBlockStatus(userId int, blocked bool) error
	UpdateUserRole(userId int, role string) error
	GetUserStatistics(userId int) (dto.UserStats, error)
	GetSystemStatistics() (dto.SystemStats, error)
	DeleteUser(userId int) error
	IsAdmin(userId int) (bool, error)
}
type Repository struct {
	Admin
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin: postgres.NewAdminRepository(db),
	}
}
