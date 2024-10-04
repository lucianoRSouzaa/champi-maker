package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" validate:"required"`
	Name         string    `json:"name" validate:"required,min=2,max=100"`
	Email        string    `json:"email" validate:"required,email,max=100"`
	PasswordHash string    `json:"-" validate:"required"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
