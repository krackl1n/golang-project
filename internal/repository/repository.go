package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krackl1n/golang-project/internal/apperr"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/pkg/errors"
)

type userRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) UserProvider {
	return &userRepository{
		conn: conn,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users(id, name, age, gender, email) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.conn.Exec(ctx, query, user.ID.String(), user.Name, user.Age, user.Gender, user.Email)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "create user")
	}

	slog.Debug(fmt.Sprintf("created user: id=%s", user.ID))
	return user.ID, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, age, gender, email 
		FROM users 
		WHERE id=$1
	`

	var user models.User
	row := r.conn.QueryRow(ctx, query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Gender, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrorNotFound
		}
		return nil, errors.Wrap(err, "scan user data")
	}

	slog.Debug(fmt.Sprintf("received successfully: id=%s", id))
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET name = $1, age = $2, gender = $3, email = $4
        WHERE id = $5
    `

	result, err := r.conn.Exec(ctx, query, user.Name, user.Age, user.Gender, user.Email, user.ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("update user: id=%s", user.ID))
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrorNotFound
	}

	slog.Debug(fmt.Sprintf("updated successfully: id=%s", user.ID))
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id=$1
	`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete user: id=%s", id))
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrorNotFound
	}

	slog.Debug(fmt.Sprintf("deleted successfully: id=%s", id))
	return nil
}
