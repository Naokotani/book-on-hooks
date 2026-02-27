package repository

import "context"

func (db *Database) Ping(ctx context.Context) error {
	return db.Db.Ping(ctx)
}
