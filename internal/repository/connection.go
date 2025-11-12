package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database struct {
	Conn  *bun.DB
	Ctx   context.Context
	sqldb *sql.DB
}

func GetDatabaseConnection() (Database, error) {
	ctx := context.Background()

	sqldb := createDBConnection()

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))
	conn := Database{
		Ctx:  ctx,
		Conn: db,
	}

	return conn, nil
}

func createDBConnection() *sql.DB {
	dsn := "postgres://postgres:secret@db:5432/postgres?sslmode=disable"
	for range 10 {
		db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		if err := db.Ping(); err == nil {
			return db
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
