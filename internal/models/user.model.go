package models

type User struct {
	Id       string `db:"id" json:"id,omitempty"`
	Email    string `db:"email" json:"email,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
}

type RegisterUser struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"User!23456789"`
}
