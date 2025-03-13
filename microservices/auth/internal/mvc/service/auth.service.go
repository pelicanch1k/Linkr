package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pelicanch1k/Linkr/auth/internal/config"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository"
)

type AuthService struct {
	repo repository.Auth
	authConfig *config.AuthConfig
}

func NewAuthService(repo repository.Auth, authConfig *config.AuthConfig) *AuthService {
	return &AuthService{repo: repo, authConfig: authConfig}
}

func (s *AuthService) CreateUser(userDTO dto.RegisterUser) (int, error) {
	userDTO.Password = s.generatePasswordHash(userDTO.Password)
	return s.repo.CreateUser(userDTO)
}

// Структура для хранения ключей и данных токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateJWT(userDTO dto.AuthorizationUser) (string, error) {
	user, err := s.repo.GetUserId(userDTO.Username, s.generatePasswordHash(userDTO.Password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.User_Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Токен действителен 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(s.authConfig.SigningKey))
}

// Десериализация JWT
func (s *AuthService) ParseJWT(tokenString string) (int, error) {
	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.authConfig.SigningKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Проверка токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.authConfig.Salt)))
}