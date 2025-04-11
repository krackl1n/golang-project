package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/krackl1n/golang-project/internal/repository"
	"github.com/pkg/errors"
)

type userUC struct {
	userRepository repository.UserProvider
}

func New(userRepository repository.UserProvider) UserProvider {
	return &userUC{
		userRepository: userRepository,
	}
}

func (uc *userUC) CreateUser(ctx context.Context, createUserDTO *models.CreateUserDTO) (uuid.UUID, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		return userId, errors.Wrap(err, "generate UUID")
	}

	user := &models.User{
		ID:     userId,
		Name:   createUserDTO.Name,
		Age:    createUserDTO.Age,
		Gender: createUserDTO.Gender,
		Email:  createUserDTO.Email,
	}

	if _, err := uc.userRepository.Create(ctx, user); err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (uc *userUC) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := uc.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "userUC get by Id")
	}

	return user, nil
}

func (uc *userUC) UpdateUser(ctx context.Context, user *models.User) error {
	if err := uc.userRepository.Update(ctx, user); err != nil {
		return errors.Wrap(err, "update user")
	}

	return nil
}

func (uc *userUC) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := uc.userRepository.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "delete user")
	}

	return nil
}
