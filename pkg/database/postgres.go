package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB() error {
	// dsn будет заполнять из env, а не из config для упрощения
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	var err error

	Pool, err = pgxpool.New(context.Background(), dsn)

	if err != nil {
		return errors.New("ConnectDB: error create pool")
	}

	return nil
}

func DisconnectDB() {
	if Pool != nil {
		Pool.Close()
	}
}
