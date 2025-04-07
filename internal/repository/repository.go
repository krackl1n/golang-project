package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

	_, err := r.conn.Exec(ctx, query, user.Id.String(), user.Name, user.Age, user.Gender, user.Email)
	if err != nil {
		slog.Error("Failed to create user", "email", user.Email, "error", err)
		return uuid.Nil, errors.Wrap(err, "create user")
	}

	slog.Info("Created user", "id", user.Id)
	return user.Id, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, age, gender, email 
		FROM users 
		WHERE id=$1
	`

	var user models.User
	row := r.conn.QueryRow(ctx, query, id)
	err := row.Scan(&user.Id, &user.Name, &user.Age, &user.Gender, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("User not found", "id", id)
			return nil, errors.New("user not found")
		}
		slog.Error("Failed to get user by ID", "id", id, "error", err)
		return nil, errors.Wrap(err, "scan user data")
	}

	slog.Info("User received successfully", "id", id)
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET name = $1, age = $2, gender = $3, email = $4
        WHERE id = $5
    `

	result, err := r.conn.Exec(ctx, query, user.Name, user.Age, user.Gender, user.Email, user.Id)
	if err != nil {
		slog.Error("Failed to update user", "id", user.Id, "error", err)
		return errors.Wrap(err, "update user")
	}

	rows := result.RowsAffected()
	if rows == 0 {
		slog.Warn("User not found for update", "id", user.Id)
		return errors.New("user not found")
	}

	slog.Info("User updated successfully", "id", user.Id)
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id=$1
	`

	result, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete user", "id", id, "error", err)
		return errors.Wrap(err, "delete user")
	}

	rows := result.RowsAffected()
	if rows == 0 {
		slog.Warn("User not found", "id", id)
		return errors.New("user not found")
	}

	slog.Info("User deleted successfully", "id", id)
	return nil
}
