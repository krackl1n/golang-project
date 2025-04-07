package usecase

import (
	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/krackl1n/golang-project/internal/repository"
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

func (uc *UseCase) CreateUser(createUserDTO *models.CreateUserDTO) (uuid.UUID, error) {

	uuidUser, err := uuid.NewV7()
	if err != nil {
		return uuidUser, err
	}
	// user := models.User
	return uc.userRepository.Create(&models.User{
		Id:     uuidUser,
		Name:   createUserDTO.Name,
		Age:    createUserDTO.Age,
		Gender: createUserDTO.Gender,
		Email:  createUserDTO.Email,
	})
}

func (uc *UseCase) GetUserById(id uuid.UUID) (*models.User, error) {
	return uc.userRepository.GetByID(id)
}

func (uc *UseCase) UpdateUser(user *models.User) error {
	return uc.userRepository.Update(user)
}

func (uc *UseCase) DeleteUser(id uuid.UUID) error {
	return uc.userRepository.Delete(id)
}
