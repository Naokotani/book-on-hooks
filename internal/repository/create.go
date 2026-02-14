package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/sqlc"
)

func (db *Database) InsertBook(ctx context.Context,
	book *domain.Book, file multipart.File,
	header *multipart.FileHeader) (int64, error) {

	tx, err := db.Db.Begin(ctx)

	if err != nil {
		return 0, fmt.Errorf("Failed to open database transaction:\n%s", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Q.WithTx(tx)

	id, err := qtx.InsertBook(ctx,
		sqlc.InsertBookParams{
			Title:   book.Title,
			Author:  book.Author,
			Summary: book.Summary,
			Price:   book.Price,
		})

	if err != nil {
		return 0, fmt.Errorf("failed to insert book\n%s", err)
	}

	filename := strconv.FormatInt(id, 10) + "_" + header.Filename

	dst, err := os.Create(filepath.Join("./images/covers", filename))
	if err != nil {
		return 0, fmt.Errorf("failed to create image file\n%s", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return 0, fmt.Errorf("failed to copy file to image directory\n%s", err)
	}

	err = qtx.UpdateBookImage(ctx, sqlc.UpdateBookImageParams{
		ID:    id,
		Image: filename,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to update book with image filename\n%s", err)
	}

	tx.Commit(ctx)
	return id, nil
}

func (db *Database) InsertMachine(ctx context.Context, machine *domain.Machine) (int64, error) {
	id, err := db.Q.InsertMachine(ctx, sqlc.InsertMachineParams{
		Location: machine.Location,
		Rows:     int32(machine.Rows),
		Columns:  int32(machine.Columns),
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
