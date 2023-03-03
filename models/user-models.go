package models

import (
	"database/sql"
	"server/store"
)

type User struct {
	*sql.DB
	Username   string `json:"username", omitempty`
	Email      string `json:"email", omitempty`
	Password   string `json:"password"`
	Repassword string `json:"repassword", omitempty`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}
type IUser interface {
	store.IDbManipulation[User]
}

func NewUser(db *sql.DB) IUser {
	return store.NewDbManipulation[User](db)
}

type Response struct {
	Code    int
	Message []string
}

type SendingJson struct {
	Bar string
}
