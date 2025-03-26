package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/profile/internal/model"
	"github.com/pelicanch1k/Linkr/profile/internal/mvc/repository/postgres"
)

type Common interface {
	GetUserById(user_id int) model.User_Common
	GetUserByUsername(username string) model.User_Common
}

type Owner interface {
	GetUser(token string) model.User_Owner

	SetUsername(user model.User_Owner, username string) model.User_Owner
	SetEmail(user model.User_Owner, email string) model.User_Owner
	SetPassword(user model.User_Owner, password_hash string) model.User_Owner
}

type Repository struct {
	Common
	Owner
}

func NewRepository(db *sqlx.DB)*Repository {
	return &Repository{
		Common: postgres.NewCommonRepository(db),
		Owner: postgres.NewOwnerRepository(db),
	}
}