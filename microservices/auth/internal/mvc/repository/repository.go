package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/model"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository/postgres"
)

type Auth interface {
	CreateUser(user dto.RegisterUser) (int, error)
	GetUserId(username, password_hash string) (model.User, error)
}

type Repository struct {
	Auth
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:     postgres.NewAuthRepository(db),
	}
}