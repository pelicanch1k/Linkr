package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/dto"
	"github.com/pelicanch1k/Linkr/auth/internal/model"

	_ "github.com/lib/pq"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) CreateUser(user dto.RegisterUser) (int, error) {
	query := "INSERT INTO users (username, email, password_hash, first_name, last_name) values ($1, $2, $3, $4, $5) RETURNING user_id"
	row := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.First_name, user.Last_name)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthRepository) GetUserId(username, password_hash string) (model.User, error) {
	var user model.User
	query := "SELECT user_id FROM users WHERE username = $1 AND password_hash = $2 LIMIT 1"
	err := r.db.Get(&user, query, username, password_hash)

	return user, err
}
