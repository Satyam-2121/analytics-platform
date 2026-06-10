package database

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
	CH       clickhouse.Conn
}