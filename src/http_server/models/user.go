package models

type User struct {
	Username string `json:"username"`
}

type UserInDB struct {
	User
	Password string `json:"password"`
}

type UserData struct {
	Username string `json:"username"`
	HashedPassword string `json:"password"`
	ID string `json:"id"`
}