package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type ConfigDatabase struct {
	PostgresPort     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	PostgresHost     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresDB       string `yaml:"db" env:"POSTGRES_DB" env-default:"postgres"`
	PostgresUser     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

type Config struct {
	ServicePort string `yaml:"service_port" env:"SERVICE_PORT" env-default:"8080"`
	MetricsPort string `yaml:"metrics_port" env:"METRICS_PORT" env-default:"8081"`
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" env-default:"INFO"`
	ConnString  string
	ConfigDatabase
}

func Load() (*Config, error) {

	var cfg Config
	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, errors.Wrap(err, "load config")
	}

	cfg.ConnString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)

	return &cfg, nil
}
