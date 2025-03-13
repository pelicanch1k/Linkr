package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/auth/internal/config/db"
)

func NewPostgresDriver(cfg db.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}