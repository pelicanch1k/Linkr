package service

import (
	"fmt"

	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository"
)

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetUsers() ([]dto.UserProfile, error) {
	return s.repo.GetAllUsers()
}

func (s *AdminService) GetUserById(userId int) (dto.UserProfile, error) {
	return s.repo.GetUserById(userId)
}

func (s *AdminService) BlockUser(userId int, blocked bool) error {
	return s.repo.UpdateUserBlockStatus(userId, blocked)
}

func (s *AdminService) ChangeUserRole(userId int, role string) error {
	// Проверка допустимости роли
	if role != "user" && role != "admin" && role != "moderator" {
		return fmt.Errorf("недопустимая роль: %s", role)
	}

	return s.repo.UpdateUserRole(userId, role)
}

func (s *AdminService) GetUserStats(userId int) (dto.UserStats, error) {
	return s.repo.GetUserStatistics(userId)
}

func (s *AdminService) GetSystemStats() (dto.SystemStats, error) {
	return s.repo.GetSystemStatistics()
}

func (s *AdminService) DeleteUser(userId int) error {
	return s.repo.DeleteUser(userId)
}
