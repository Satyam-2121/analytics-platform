package database

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func ConnectClickHouse() (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "admin",
			Password: "admin",
		
		},
		
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	fmt.Println("Connected to ClickHouse")

	return conn, nil
}