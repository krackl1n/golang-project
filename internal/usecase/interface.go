package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
)

type UserProvider interface {
	CreateUser(ctx context.Context, createUserDTO *models.CreateUserDTO) (uuid.UUID, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
