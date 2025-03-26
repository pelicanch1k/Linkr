package model

type User_Common struct{
	User_id int
	Username string
	FirstName string
	LastName string
	Bio string
	ProfilePictureUrl string
}

type User_Owner struct{
	User_Common
	Email string
}