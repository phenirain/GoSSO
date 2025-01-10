package models

type User struct {
	Id           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email" binding:"required"`
	PasswordHash []byte `json:"password_hash" db:"password_hash"`
}
