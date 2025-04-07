package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
)

type UserProvider interface {
	Create(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
