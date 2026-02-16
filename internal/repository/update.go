package repository

import (
	"context"

	"booksonhooks.ca/internal/domain"
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

func (db *Database) UpdateBook(ctx context.Context, book *domain.Book) error {
	err := db.Q.UpdateBook(ctx, sqlc.UpdateBookParams{
		ID:      book.ID,
		Title:   book.Title,
		Author:  book.Author,
		Summary: book.Summary,
		Price:   book.Price,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteBook(ctx context.Context, id int64) error {
	err := db.Q.DeleteBook(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateMachine(ctx context.Context, machine *domain.Machine) error {
	err := db.Q.UpdateMachine(ctx, sqlc.UpdateMachineParams{
		ID:       machine.ID,
		Location: machine.Location,
		Rows:     int32(machine.Rows),
		Columns:  int32(machine.Columns),
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteMachine(ctx context.Context, id int64) error {
	err := db.Q.DeleteMachine(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ClearMachineBooks(ctx context.Context, machineID int64) error {
	err := db.Q.DeleteBookMachineByMachineID(ctx, machineID)
	if err != nil {
		return err
	}

	return nil
}
