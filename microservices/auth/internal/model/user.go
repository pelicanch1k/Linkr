package model

type User struct {
	User_Id int `json:"-" db:"user_id"`
}
