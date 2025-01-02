package conn

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Pool() *pgxpool.Pool {

	s, b := os.LookupEnv("GOCRIB_HOST")

	host := s
	if !b {
		host = "db"
	}

	dsn := "postgres://postgres:example@" + host + ":5432/cribbage?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
