package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type DB struct {
	pool *pgxpool.Pool
}

func connectToDB(c Config) (*pgxpool.Pool, error) {
	user := os.Getenv(c.UserEnvKey)
	password := os.Getenv(c.PassEnvKey)
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", c.Addr, c.Port, c.DB, user, password)

	return pgxpool.Connect(context.Background(), connectionString)
}

func NewDB(c Config) (*DB, error) {
	db, err := connectToDB(c)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to db")
	}

	return &DB{pool: db}, nil
}
