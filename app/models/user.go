package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `json:"username" gorm:"unique;type:varchar(20)" validate:"required,min=6,max=32"`
	Password  string `json:"password" gorm:"type:varchar(255)" validate:"required,min=6"`
	Fullname  string `json:"full_name" gorm:"type:varchar(100)" validate:"required,min=6"`
}

func (l User) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserSession struct {
	ID                  uint `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UserID              uint `json:"user_id" gorm:"type:int" validate:"required"`
	Token               string `json:"token" gorm:"type:varchar(255)" validate:"required"`
	RefreshToken        string `json:"refresh_token" gorm:"type:varchar(255)" validate:"required"`
	TokenExpired time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
type LoginResponse struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}