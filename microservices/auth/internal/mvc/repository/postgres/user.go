package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"

	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserProfile(userId int) (dto.UserProfile, error) {
	var profile dto.UserProfile
	query := `SELECT user_id as id, username, email, first_name, last_name, profile_picture_url, role, created_at
              FROM users 
              WHERE user_id = $1`

	err := r.db.Get(&profile, query, userId)
	if err != nil {
		return dto.UserProfile{}, err
	}

	return profile, nil
}

func (r *UserRepository) UpdateUserProfile(userId int, profile dto.UpdateProfileRequest) error {
	query := `UPDATE users SET 
              first_name = COALESCE(NULLIF($1, ''), first_name),
              last_name = COALESCE(NULLIF($2, ''), last_name),
              username = COALESCE(NULLIF($3, ''), username),
              updated_at = CURRENT_TIMESTAMP
              WHERE user_id = $4 AND (deleted IS NULL OR deleted = false)`

	_, err := r.db.Exec(query, profile.FirstName, profile.LastName, profile.Username, userId)
	return err
}

func (r *UserRepository) DeleteUser(userId int) error {
	// Используем мягкое удаление
	query := `UPDATE users SET 
              deleted = true,
              deleted_at = CURRENT_TIMESTAMP
              WHERE user_id = $1`

	_, err := r.db.Exec(query, userId)
	return err
}

func (r *UserRepository) UpdateAvatar(userId int, avatarURL string) error {
	query := `UPDATE users SET 
              profile_picture_url = $1,
              updated_at = CURRENT_TIMESTAMP
              WHERE user_id = $2 AND (deleted IS NULL OR deleted = false)`

	_, err := r.db.Exec(query, avatarURL, userId)
	return err
}

func (r *UserRepository) UpdateEmail(userId int, email string) error {
	query := `UPDATE users SET 
              email = $1,
              updated_at = CURRENT_TIMESTAMP
              WHERE user_id = $2 AND (deleted IS NULL OR deleted = false)`

	_, err := r.db.Exec(query, email, userId)
	return err
}

func (r *UserRepository) CheckUsernameExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users 
              WHERE username = $1 AND (deleted IS NULL OR deleted = false)`

	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) CheckEmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users 
              WHERE email = $1 AND (deleted IS NULL OR deleted = false)`

	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) GetUserNotifications(userId int) ([]dto.Notification, error) {
	var notifications []dto.Notification
	query := `SELECT notification_id as id, user_id, message, read, created_at
              FROM user_notifications
              WHERE user_id = $1
              ORDER BY created_at DESC
              LIMIT 50`

	err := r.db.Select(&notifications, query, userId)
	if err != nil {
		// Если нет записей, это не ошибка
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.Notification{}, nil
		}
		return nil, err
	}

	return notifications, nil
}
