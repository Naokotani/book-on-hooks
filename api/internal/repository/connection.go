package repository

import (
	"context"
	"os"
	"time"

	"booksonhooks.ca/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Q  *sqlc.Queries
	Db *pgxpool.Pool
}

func GetDatabaseConnection() (*Database, error) {
	ctx := context.Background()

	db, err := createDBConnection(ctx)

	if err != nil {
		return nil, err
	}

	q := sqlc.New(db)

	database := Database{
		Db: db,
		Q:  q,
	}

	return &database, nil
}

func createDBConnection(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://postgres:secret@db:5432/books?sslmode=disable"
	}
	var err error
	for range 10 {
		conn, err := pgxpool.New(ctx, url)
		if err == nil {
			if pingErr := conn.Ping(ctx); pingErr == nil {
				return conn, nil
			} else {
				err = pingErr
				conn.Close()
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil, err
}
