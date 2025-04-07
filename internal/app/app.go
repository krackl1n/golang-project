package app

import (
	"fmt"
	"log/slog"

	"github.com/krackl1n/golang-project/database"
	"github.com/krackl1n/golang-project/internal/config"
	"github.com/krackl1n/golang-project/internal/handler"
	"github.com/krackl1n/golang-project/internal/repository"
	"github.com/krackl1n/golang-project/internal/storage"
	"github.com/krackl1n/golang-project/internal/usecase"
	"github.com/pkg/errors"
)

func Run() error {
	// Загрузка конфигов
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		return errors.Wrap(err, "failed to load config")
	}

	// Миграции
	if err := database.Migrate(cfg.ConnString); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		return errors.Wrap(err, "failed to apply migrations")
	}

	// Коннект к бд
	connDB, err := storage.GetConnect(cfg.ConnString)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return errors.Wrap(err, "failed to connect to database")
	}
	defer connDB.Close()

	// DI
	userRepository := repository.NewUserRepository(connDB)
	usecase := usecase.NewUseCase(userRepository)
	handler := handler.NewHandler(usecase)

	// Получение роутера
	app := GetRouter(handler)

	if err := app.Listen(fmt.Sprintf(":%s", cfg.ServicePort)); err != nil {
		slog.Error("Application listening error", "error", err)
		return errors.Wrap(err, "application listening error")
	}
	return nil
}
