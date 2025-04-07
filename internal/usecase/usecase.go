package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/krackl1n/golang-project/internal/repository"
	"github.com/pkg/errors"
)

type UseCase struct {
	userRepository repository.UserProvider
}

func NewUseCase(userRepository repository.UserProvider) *UseCase {
	return &UseCase{
		userRepository: userRepository,
	}
}

// TODO Дописать usecase

func (uc *UseCase) CreateUser(ctx context.Context, createUserDTO *models.CreateUserDTO) (uuid.UUID, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return userId, errors.Wrap(err, "generate UUID")
	}

	user := &models.User{
		Id:     userId,
		Name:   createUserDTO.Name,
		Age:    createUserDTO.Age,
		Gender: createUserDTO.Gender,
		Email:  createUserDTO.Email,
	}

	if _, err := uc.userRepository.Create(ctx, user); err != nil {
		slog.Error("Failed to create user", "user_id", userId, "error", err)
		return uuid.Nil, errors.Wrap(err, "usecase create user")
	}

	slog.Info("User created successfully", "user_id", userId)
	return userId, nil
}

func (uc *UseCase) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := uc.userRepository.GetByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user by Id", "id", id, "error", err)
		return nil, errors.Wrap(err, "usecase get user by Id")
	}

	slog.Info("User received successfully", "id", id, "error", err)
	return user, nil
}

func (uc *UseCase) UpdateUser(ctx context.Context, user *models.User) error {
	if err := uc.userRepository.Update(ctx, user); err != nil {
		slog.Error("Failed to update user", "id", user.Id, "error", err)
		return errors.Wrap(err, "update user")
	}

	slog.Info("User updated successfully", "id", user.Id)
	return nil
}

func (uc *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := uc.userRepository.Delete(ctx, id); err != nil {
		slog.Error("Failed to delete user", "id", id, "error", err)
		return errors.Wrap(err, "delete user")
	}

	slog.Info("User deleted successfully", "id", id)
	return nil
}
