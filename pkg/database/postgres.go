package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// глобальная перменная структуры пула подключения к бд
var Pool *pgxpool.Pool

func ConnectDB() error {
	// dsn будет заполняться из env, а не из config для упрощения
	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		return errors.New("ConnectDB: DB_URL is not set")
	}

	// Создаем пул
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("ConnectDB: error create pool: %w", err)
	}

	// Проверяем соединение
	err = Pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("ConnectDB: database is unreachable: %w", err)
	}

	return nil
}

func DisconnectDB() {
	if Pool != nil {
		Pool.Close()
	}
}
