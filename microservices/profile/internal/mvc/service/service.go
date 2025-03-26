package service

import "github.com/pelicanch1k/Linkr/profile/internal/model"

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