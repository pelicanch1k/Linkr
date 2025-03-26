package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/Linkr/profile/internal/model"
)

type OwnerRepository struct {
	db *sqlx.DB
}

func NewOwnerRepository(db *sqlx.DB) *OwnerRepository {
	return &OwnerRepository{db}
}

func (o *OwnerRepository) GetUser(token string) model.User_Owner{
	
}

func (o *OwnerRepository) SetUsername(user model.User_Owner, username string) model.User_Owner {

}

func (o *OwnerRepository) SetEmail(user model.User_Owner, email string) model.User_Owner {

}

func (o *OwnerRepository) SetPassword(user model.User_Owner, password_hash string) model.User_Owner {

}