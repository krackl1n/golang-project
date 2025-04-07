package config

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	ServicePort      string `mapstructure:"SERVICE_PORT"`
	ConnString       string
}

func Load() (*Config, error) {
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("SERVICE_PORT", "8080")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("Config file not found, using environment variables only")
		} else {
			slog.Error("Error reading config file", "error", err)
			return nil, errors.Wrap(err, "failed to read config")
		}
	}

	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		slog.Error("Error unmarshalling config", "error", err)
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	if err := cfg.validate(); err != nil {
		slog.Error("Configuration validation failed", "error", err)
		return nil, errors.Wrap(err, "invalid configuration")
	}

	cfg.ConnString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)

	return cfg, nil
}

func (c *Config) validate() error {
	if c.PostgresUser == "" {
		return errors.New("POSTGRES_USER is required")
	}
	if c.PostgresPassword == "" {
		return errors.New("POSTGRES_PASSWORD is required")
	}
	if c.PostgresHost == "" {
		return errors.New("POSTGRES_HOST is required")
	}
	if c.PostgresDB == "" {
		return errors.New("POSTGRES_DB is required")
	}
	if c.ServicePort == "" {
		return errors.New("SERVICE_PORT is required")
	}
	return nil
}
