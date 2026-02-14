package repository

import (
	"context"

	"booksonhooks.ca/internal/sqlc"
)

func (db *Database) UpdateBookPosition(ctx context.Context, id int64, row, col int32) error {
	err := db.Q.UpdateBookPosition(ctx, sqlc.UpdateBookPositionParams{
		BookID: id,
		Row:    row,
		Col:    col,
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateCoverFilename(ctx context.Context, id int64, image string) error {
	err := db.Q.UpdateBookImage(ctx, sqlc.UpdateBookImageParams{
		ID:    id,
		Image: image,
	})

	if err != nil {
		return err
	}

	return nil
}
