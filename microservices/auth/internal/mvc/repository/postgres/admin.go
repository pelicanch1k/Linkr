package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"

	_ "github.com/lib/pq"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{db}
}

func (r *AdminRepository) GetAllUsers() ([]dto.UserProfile, error) {
	var users []dto.UserProfile
	query := `SELECT user_id as id, username, email, first_name, last_name, 
                     profile_picture_url as avatar, role, created_at
              FROM users
              WHERE (deleted IS NULL OR deleted = false)
              ORDER BY created_at DESC`

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *AdminRepository) GetUserById(userId int) (dto.UserProfile, error) {
	var user dto.UserProfile
	query := `SELECT user_id as id, username, email, first_name, last_name, 
                     profile_picture_url as avatar, role, created_at
              FROM users
              WHERE user_id = $1 AND (deleted IS NULL OR deleted = false)`

	err := r.db.Get(&user, query, userId)
	if err != nil {
		return dto.UserProfile{}, err
	}

	return user, nil
}

func (r *AdminRepository) UpdateUserBlockStatus(userId int, blocked bool) error {
	query := `UPDATE users SET 
              blocked = $1,
              updated_at = CURRENT_TIMESTAMP
              WHERE user_id = $2`

	_, err := r.db.Exec(query, blocked, userId)
	return err
}

func (r *AdminRepository) UpdateUserRole(userId int, role string) error {
	query := `UPDATE users SET 
              role = $1,
              updated_at = CURRENT_TIMESTAMP
              WHERE user_id = $2`

	_, err := r.db.Exec(query, role, userId)
	return err
}

func (r *AdminRepository) GetUserStatistics(userId int) (dto.UserStats, error) {
	var stats dto.UserStats
	query := `SELECT last_login, COALESCE(login_count, 0) as login_count, 
                    (SELECT COUNT(*) FROM user_tokens WHERE user_id = users.user_id) as session_count
              FROM users
              WHERE user_id = $1`

	err := r.db.Get(&stats, query, userId)
	if err != nil {
		return dto.UserStats{}, err
	}

	return stats, nil
}

func (r *AdminRepository) GetSystemStatistics() (dto.SystemStats, error) {
	var stats dto.SystemStats

	// Получаем общее количество пользователей
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE (deleted IS NULL OR deleted = false)`).
		Scan(&stats.UserCount)
	if err != nil {
		return dto.SystemStats{}, err
	}

	// Получаем активных пользователей (логин за последний месяц)
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users 
                         WHERE last_login > CURRENT_TIMESTAMP - INTERVAL '30 days' 
                         AND (deleted IS NULL OR deleted = false)`).
		Scan(&stats.ActiveUsers)
	if err != nil {
		return dto.SystemStats{}, err
	}

	// Получаем заблокированных пользователей
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users 
                         WHERE blocked = true AND (deleted IS NULL OR deleted = false)`).
		Scan(&stats.LockedAccounts)
	if err != nil {
		return dto.SystemStats{}, err
	}

	// Получаем удаленных пользователей
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE deleted = true`).
		Scan(&stats.DeletedAccounts)
	if err != nil {
		return dto.SystemStats{}, err
	}

	return stats, nil
}

func (r *AdminRepository) DeleteUser(userId int) error {
	// Используем мягкое удаление
	query := `UPDATE users SET 
              deleted = true,
              deleted_at = CURRENT_TIMESTAMP
              WHERE user_id = $1`

	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}

	// Удаляем токены пользователя
	_, err = r.db.Exec(`DELETE FROM user_tokens WHERE user_id = $1`, userId)
	return err
}
