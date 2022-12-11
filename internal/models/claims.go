package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Claims struct {
	ID        int64     `json:"id" validate:"required"`
	Role      Role      `json:"role" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
}

func (c *Claims) Valid() error {
	validate := validator.New()
	return validate.Struct(c)
}
