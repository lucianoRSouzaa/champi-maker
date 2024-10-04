package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Logo      string    `json:"logo" validate:"omitempty,url"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func (t *Team) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
