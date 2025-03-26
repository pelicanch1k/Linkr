package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/profile/internal/model"
)

type CommonRepository struct {
	db *sqlx.DB
}

func NewCommonRepository(db *sqlx.DB) *CommonRepository {
	return &CommonRepository{db}
}

func (c *CommonRepository) GetUserById(user_id int) model.User_Common{
	return model.User_Common{}
}

func (c *CommonRepository) GetUserByUsername(username string) model.User_Common {
	return model.User_Common{}
}