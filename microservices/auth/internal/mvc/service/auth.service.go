package service

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pelicanch1k/Linkr/auth/internal/config"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository"
)

// Claims для JWT токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo       repository.Auth
	authConfig *config.AuthConfig
}

func NewAuthService(repo repository.Auth, authConfig *config.AuthConfig) *AuthService {
	return &AuthService{repo: repo, authConfig: authConfig}
}

func (s *AuthService) CreateUser(user dto.RegisterUser) (int, error) {
	
	// Валидация email
	if !s.isValidEmail(user.Email) {
		return 0, fmt.Errorf("некорректный формат email")
	}

	// Проверка имени и фамилии
	if strings.TrimSpace(user.First_name) == "" {
		return 0, fmt.Errorf("имя не может быть пустым")
	}

	if strings.TrimSpace(user.Last_name) == "" {
		return 0, fmt.Errorf("фамилия не может быть пустой")
	}

	// Проверка наличия email в базе
	exists, err := s.repo.CheckEmailExists(user.Email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("пользователь с таким email уже существует")
	}

	// Проверка наличия username в базе
	exists, err = s.repo.CheckUsernameExists(user.Username)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("пользователь с таким username уже существует")
	}

	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateJWT(user dto.AuthorizationUser) (string, error) {
	user.Password = s.generatePasswordHash(user.Password)

	userId, err := s.repo.GetUserId(user)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(s.authConfig.SigningKey))
	if err != nil {
		return "", err
	}

	// Сохраняем токен в базе
	if err := s.repo.StoreToken(userId, tokenString); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseJWT(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
		}
		return []byte(s.authConfig.SigningKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("недействительный токен")
	}

	// Проверка, не был ли токен отозван
	isValid, err := s.repo.CheckToken(claims.UserID, tokenString)
	if err != nil {
		return 0, err
	}

	if !isValid {
		return 0, fmt.Errorf("токен отозван")
	}

	return claims.UserID, nil
}

func (s *AuthService) Logout(userId int) error {
	return s.repo.InvalidateTokens(userId)
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	userId, err := s.ParseJWT(refreshToken)
	if err != nil {
		return "", fmt.Errorf("недействительный токен обновления: %w", err)
	}

	// Инвалидируем старый токен
	if err := s.repo.InvalidateToken(userId, refreshToken); err != nil {
		return "", err
	}

	// Генерируем новый токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(s.authConfig.SigningKey))
	if err != nil {
		return "", err
	}

	// Сохраняем новый токен
	if err := s.repo.StoreToken(userId, tokenString); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) SendResetPasswordEmail(email string) error {
	// Проверка существования email
	exists, err := s.repo.CheckEmailExists(email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("пользователь с таким email не найден")
	}

	// Создание токена сброса пароля
	token := fmt.Sprintf("%d", time.Now().Unix())
	hashToken := s.generatePasswordHash(token)

	// Сохранение токена в базе данных
	if err := s.repo.StoreResetToken(email, hashToken); err != nil {
		return err
	}

	// Здесь должна быть логика отправки email с токеном
	// В данной реализации просто возвращаем успех
	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Проверка токена
	userId, err := s.repo.ValidateResetToken(token)
	if err != nil {
		return fmt.Errorf("недействительный токен сброса: %w", err)
	}

	// Хеширование нового пароля
	hashedPassword := s.generatePasswordHash(newPassword)

	// Обновление пароля
	return s.repo.UpdatePassword(userId, hashedPassword)
}

func (s *AuthService) ChangePassword(userId int, oldPassword, newPassword string) error {
	// Проверка старого пароля
	hashedOldPassword := s.generatePasswordHash(oldPassword)
	isValid, err := s.repo.CheckPassword(userId, hashedOldPassword)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("неверный текущий пароль")
	}

	// Хеширование и обновление пароля
	hashedNewPassword := s.generatePasswordHash(newPassword)
	return s.repo.UpdatePassword(userId, hashedNewPassword)
}

func (s *AuthService) IsAdmin(userId int) (bool, error) {
	return s.repo.IsAdmin(userId)
}

func (s *AuthService) ValidateToken(token string) (bool, error) {
	_, err := s.ParseJWT(token)
	return err == nil, nil
}

func (s *AuthService) RevokeAllTokens(userId int) error {
	return s.repo.InvalidateTokens(userId)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Write([]byte(s.authConfig.Salt))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Проверка формата email с помощью регулярного выражения
func (s *AuthService) isValidEmail(email string) bool {
	// Простая проверка регулярным выражением
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
