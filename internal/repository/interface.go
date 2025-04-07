package repository

import (
	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
)

type UserProvider interface {
	Create(user *models.User) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}
