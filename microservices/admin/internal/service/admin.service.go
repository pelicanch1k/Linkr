package service

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/pelicanch1k/Linkr/admin/internal/dto"
	"github.com/pelicanch1k/Linkr/admin/internal/repository"
)

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetUsers() ([]dto.UserProfile, error) {
	usersModel, err := s.repo.GetAllUsers()
	if err != nil {
		return []dto.UserProfile{}, err
	}

	var usersDTO []dto.UserProfile

	err = copier.Copy(&usersDTO, &usersModel)
	if err != nil {
		return []dto.UserProfile{}, err
	}

	return usersDTO, nil
}

func (s *AdminService) GetUserById(userId int) (dto.UserProfile, error) {
	userModel, err := s.repo.GetUserById(userId)
	if err != nil {
		return dto.UserProfile{}, err
	}

	var userDTO dto.UserProfile

	err = copier.Copy(&userDTO, &userModel)
	if err != nil {
		return dto.UserProfile{}, err
	}

	return userDTO, nil
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

func (s *AdminService) IsAdmin(userId int) (bool, error) {
	return s.repo.IsAdmin(userId)
}