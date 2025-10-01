package models

import (
	"mime/multipart"
	"time"
)

type UserProfile struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email,omitempty"`
	Name      string    `json:"name,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditUserProfile struct {
	Name   string                `form:"name"`
	Bio    string                `form:"bio"`
	Avatar *multipart.FileHeader `form:"avatar"`
}
