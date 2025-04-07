package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnect(connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connString)
}
