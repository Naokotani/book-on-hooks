package repository

import (
	"context"
	"time"

	"booksonhooks.ca/internal/sqlc"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	Q  *sqlc.Queries
	Db *pgx.Conn
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

func createDBConnection(ctx context.Context) (*pgx.Conn, error) {
	url := "postgres://postgres:secret@db:5432/postgres?sslmode=disable"
	var err error
	for range 10 {
		conn, err := pgx.Connect(ctx, url)
		if err == nil {
			return conn, nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil, err
}
