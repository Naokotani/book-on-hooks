package repository

import (
	"context"

	"booksonhooks.ca/internal/sqlc"
)

func (db *Database) MachineHasLoadedBooksOutsideBounds(ctx context.Context, machineID int64, rows int, cols int) (bool, error) {
	return db.Q.MachineHasLoadedBooksOutsideBounds(ctx, sqlc.MachineHasLoadedBooksOutsideBoundsParams{
		MachineID: machineID,
		Rows:      int32(rows),
		Cols:      int32(cols),
	})
}
