package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(userId int) (dto.UserProfile, error) {
	return s.repo.GetUserProfile(userId)
}

func (s *UserService) UpdateProfile(userId int, profile dto.UpdateProfileRequest) error {
	// Проверка на уникальность username
	if profile.Username != "" {
		exists, err := s.repo.CheckUsernameExists(profile.Username)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("этот username уже используется")
		}
	}

	return s.repo.UpdateUserProfile(userId, profile)
}

func (s *UserService) DeleteProfile(userId int) error {
	return s.repo.DeleteUser(userId)
}

func (s *UserService) UpdateAvatar(userId int, file *multipart.FileHeader) (string, error) {
	// Проверяем расширение файла
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("недопустимый формат файла: поддерживаются только .jpg, .jpeg и .png")
	}

	// Генерируем имя файла
	filename := fmt.Sprintf("%d_%d%s", userId, time.Now().Unix(), ext)

	// Сохраняем файл и пути в БД
	// Реальная реализация будет зависеть от способа хранения файлов
	avatarURL := fmt.Sprintf("/avatars/%s", filename)

	if err := s.repo.UpdateAvatar(userId, avatarURL); err != nil {
		return "", err
	}

	return avatarURL, nil
}

func (s *UserService) UpdateEmail(userId int, email string) error {
	// Проверка на уникальность email
	exists, err := s.repo.CheckEmailExists(email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("этот email уже используется")
	}

	return s.repo.UpdateEmail(userId, email)
}

func (s *UserService) GetNotifications(userId int) ([]dto.Notification, error) {
	return s.repo.GetUserNotifications(userId)
}
