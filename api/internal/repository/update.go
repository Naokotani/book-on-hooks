package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

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

func (db *Database) UpdateBookImage(ctx context.Context, id int64, file multipart.File, header *multipart.FileHeader) (string, error) {
	book, err := db.Q.GetBookByID(ctx, id)
	if err != nil {
		return "", err
	}

	if err != nil {
		_ = os.Remove(filepath.Join(imageDir(), "covers", book.Image))
		return "", err
	}

	newFilename, err := storeBookImage(id, file, header)
	if err != nil {
		return "", err
	}

	err = db.Q.UpdateBookImage(ctx, sqlc.UpdateBookImageParams{
		ID:    id,
		Image: newFilename,
	})

	return deleteBookImage(book.Image), nil
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

func (db *Database) DeleteBook(ctx context.Context, id int64) (string, error) {
	book, err := db.Q.GetBookByID(ctx, id)
	if err != nil {
		return "", err
	}

	err = db.Q.DeleteBook(ctx, id)
	if err != nil {
		return "", err
	}

	return deleteBookImage(book.Image), nil
}

func (db *Database) UpdateMachine(ctx context.Context, machine *domain.Machine) error {
	err := db.Q.UpdateMachine(ctx, sqlc.UpdateMachineParams{
		ID:       machine.ID,
		Location: machine.Location,
		Rows:     int32(machine.Rows),
		Columns:  int32(machine.Cols),
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

func deleteBookImage(filename string) string {
	if filename == "" {
		return "invariant violation: empty image filename for deleted book"
	}

	clean := filepath.Base(filename)
	if clean != filename || strings.Contains(clean, "..") {
		return fmt.Sprintf("refused to delete suspicious image filename: %q", filename)
	}

	path := filepath.Join(imageDir(), "covers", clean)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		return fmt.Sprintf("failed to delete book image %q: %v", clean, err)
	}

	return ""
}
