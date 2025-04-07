package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
)

type userRepository struct {
	users map[uuid.UUID]models.User
	mu    sync.RWMutex
}

func NewUserRepository() UserProvider {
	return &userRepository{
		users: make(map[uuid.UUID]models.User),
	}
}

func (r *userRepository) Create(user *models.User) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.Id] = *user

	return user.Id, nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[user.Id]
	if !exists {
		return fmt.Errorf("user not found")
	}

	r.users[user.Id] = *user
	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)
	return nil
}

// type userRepositoryLocal struct {
// 	conn *pgxpool.Pool
// }

// func NewRepositotyUserLocal(conn *pgxpool.Pool) UserProvider {
// 	return &userRepositoryLocal{
// 		conn: conn,
// 	}
// }

// func (r *userRepositoryLocal) Create(user *models.User) (uint64, error) {

// 	return 0, nil
// }

// func (r *userRepositoryLocal) GetByID(id uint64) (*models.User, error) {
// 	return , nil
// }

// func (r *userRepositoryLocal) Update(user *models.User) error {
// 	return nil
// }

// func (r *userRepositoryLocal) Delete(id uint64) error {
// 	return nil
// }
