package models

type User struct {
	Username string `json:"username"`
}

type UserInDB struct {
	User
	HashedPassword string `json:"hashed_password"`
}

type UserData struct {
	Username string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	ID string `json:"id"`
}