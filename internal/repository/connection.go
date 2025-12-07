package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Db *pgx.Conn
}

func GetDatabaseConnection() (*Database, error) {
	ctx := context.Background()

	db, err := createDBConnection(ctx)

	if err != nil {
		return nil, err
	}
	defer db.Close(ctx)

	database := Database{
		Db: db,
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
