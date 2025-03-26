package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"

	_ "github.com/lib/pq"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) CreateUser(user dto.RegisterUser) (int, error) {
	var id int
	query := `INSERT INTO users (username, email, password_hash, first_name, last_name, role, created_at) 
             VALUES ($1, $2, $3, $4, $5, 'user', CURRENT_TIMESTAMP) 
             RETURNING user_id`

	row := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.First_name, user.Last_name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthRepository) GetUserId(user dto.AuthorizationUser) (int, error) {
	var id int
	query := `SELECT user_id FROM users 
              WHERE username = $1 AND password_hash = $2
              LIMIT 1`

	err := r.db.QueryRow(query, user.Username, user.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("неверный логин или пароль")
		}
		return 0, err
	}

	// Обновляем дату последнего логина
	_, err = r.db.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP, login_count = COALESCE(login_count, 0) + 1 WHERE user_id = $1", id)
	if err != nil {
		// Логируем ошибку, но продолжаем выполнение
		// logger.Warn("Failed to update last_login: %v", err)
	}

	return id, nil
}

func (r *AuthRepository) CheckPassword(userId int, hashedPassword string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users 
              WHERE user_id = $1 AND password_hash = $2`

	err := r.db.QueryRow(query, userId, hashedPassword).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthRepository) UpdatePassword(userId int, hashedPassword string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2`
	_, err := r.db.Exec(query, hashedPassword, userId)
	return err
}

func (r *AuthRepository) CheckEmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND (deleted IS NULL OR deleted = false)`
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthRepository) CheckUsernameExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1 AND (deleted IS NULL OR deleted = false)`
	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthRepository) StoreToken(userId int, token string) error {
	query := `INSERT INTO user_tokens (user_id, token, created_at, expires_at) 
              VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '24 hour')`
	_, err := r.db.Exec(query, userId, token)
	return err
}

func (r *AuthRepository) CheckToken(userId int, token string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_tokens 
              WHERE user_id = $1 AND token = $2 AND expires_at > CURRENT_TIMESTAMP`
	err := r.db.QueryRow(query, userId, token).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthRepository) InvalidateToken(userId int, token string) error {
	query := `DELETE FROM user_tokens 
              WHERE user_id = $1 AND token = $2`
	_, err := r.db.Exec(query, userId, token)
	return err
}

func (r *AuthRepository) InvalidateTokens(userId int) error {
	query := `DELETE FROM user_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *AuthRepository) StoreResetToken(email string, token string) error {
	// Сначала получаем userId по email
	var userId int
	query := `SELECT user_id FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&userId)
	if err != nil {
		return err
	}

	// Удаляем предыдущие токены сброса пароля для этого пользователя
	_, err = r.db.Exec(`DELETE FROM password_reset_tokens WHERE user_id = $1`, userId)
	if err != nil {
		return err
	}

	// Создаем новый токен сброса
	query = `INSERT INTO password_reset_tokens (user_id, token, created_at, expires_at) 
             VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '24 hour')`
	_, err = r.db.Exec(query, userId, token)
	return err
}

func (r *AuthRepository) ValidateResetToken(token string) (int, error) {
	var userId int
	query := `SELECT user_id FROM password_reset_tokens 
              WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP
              LIMIT 1`
	err := r.db.QueryRow(query, token).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("недействительный или просроченный токен сброса")
		}
		return 0, err
	}

	// Удаляем использованный токен
	_, err = r.db.Exec(`DELETE FROM password_reset_tokens WHERE token = $1`, token)
	if err != nil {
		// Логируем ошибку, но продолжаем выполнение
		// logger.Warn("Failed to delete used reset token: %v", err)
	}

	return userId, nil
}

func (r *AuthRepository) IsAdmin(userId int) (bool, error) {
	var role string
	query := `SELECT role FROM users WHERE user_id = $1`
	err := r.db.QueryRow(query, userId).Scan(&role)
	if err != nil {
		return false, err
	}

	return role == "admin", nil
}
