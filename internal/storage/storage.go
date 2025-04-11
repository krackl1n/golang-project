package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConnect(connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, connString)

	if err := conn.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "ping database")
	}
	return conn, errors.Wrap(err, "get connections")
}
