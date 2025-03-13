package service

import (
	"github.com/pelicanch1k/Linkr/auth/internal/config"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository"
)

type Auth interface {
	CreateUser(user dto.RegisterUser) (int, error)
	GenerateJWT(user dto.AuthorizationUser) (string, error)
	ParseJWT(tokenString string) (int, error)
}

type Service struct {
	Auth
}

func NewService(repo repository.Auth, authConfig *config.AuthConfig) *Service {
	return &Service{
		Auth:     NewAuthService(repo, authConfig),
	}
}