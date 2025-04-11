package app

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/krackl1n/golang-project/config"
	"github.com/krackl1n/golang-project/database"
	"github.com/krackl1n/golang-project/internal/cache"
	"github.com/krackl1n/golang-project/internal/handler"
	"github.com/krackl1n/golang-project/internal/metrics"
	"github.com/krackl1n/golang-project/internal/repository"
	"github.com/krackl1n/golang-project/internal/storage"
	"github.com/krackl1n/golang-project/internal/usecase"
	"github.com/pkg/errors"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "load config")
	}
	slog.Debug("config loaded")

	loggerInit(cfg)
	slog.Debug("logger initialized")

	if err := database.Migrate(cfg.ConnString); err != nil {
		return errors.Wrap(err, "migrations")
	}
	slog.Debug("migrations applied")

	connDB, err := storage.GetConnect(cfg.ConnString)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}
	defer connDB.Close()
	slog.Debug("db connection")

	metrics.MetricsInit(cfg)
	slog.Debug("metrics initialized")

	userRepository := repository.NewUserRepository(connDB)
	userCache := cache.New(userRepository, 5*time.Minute)
	defer userCache.Stop()
	uc := usecase.New(userCache)
	handle := handler.New(uc)

	app := getRouter(handle)
	slog.Info(fmt.Sprintf("starting main server on port %s", cfg.ServicePort))
	if err := app.Listen(fmt.Sprintf(":%s", cfg.ServicePort)); err != nil {
		return errors.Wrap(err, "app listening error")
	}
	return nil
}

func loggerInit(cfg *config.Config) {
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		slog.Warn("logger init", slog.Any("error", err))
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}
