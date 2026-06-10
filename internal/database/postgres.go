package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(
		context.Background(),
		"postgres://admin:admin@localhost:5432/analytics",
	)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL")
	return pool, nil
}