package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID     uuid.UUID `json:"id" validate:"required,uuid"`
	Name   string    `json:"name" validate:"required"`
	Age    uint8     `json:"age" validate:"required,gte=0,lte=120"`
	Gender string    `json:"gender" validate:"required,oneof=male female"`
	Email  string    `json:"email" validate:"required,email"`
}

// DTO

type CreateUserDTO struct {
	Name   string `json:"name" validate:"required"`
	Age    uint8  `json:"age" validate:"required"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
	Email  string `json:"email" validate:"required,email"`
}
