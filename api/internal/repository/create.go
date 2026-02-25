package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

const defaultImageDir = "/data/images"

func imageDir() string {
	if dir := os.Getenv("IMAGE_DIR"); dir != "" {
		return dir
	}
	return defaultImageDir
}

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

	filename, err := storeBookImage(id, file, header)
	if err != nil {
		return 0, err
	}

	err = qtx.UpdateBookImage(ctx, sqlc.UpdateBookImageParams{
		ID:    id,
		Image: filename,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to update book with image filename\n%s", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction\n%s", err)
	}
	return id, nil
}

func storeBookImage(bookID int64, file multipart.File, header *multipart.FileHeader) (string, error) {
	if file == nil || header == nil {
		return "", fmt.Errorf("missing image file")
	}

	filename := strconv.FormatInt(bookID, 10) + "_" + header.Filename
	coversDir := filepath.Join(imageDir(), "covers")

	if err := os.MkdirAll(coversDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create images directory\n%s", err)
	}

	dst, err := os.Create(filepath.Join(coversDir, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create image file\n%s", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file to image directory\n%s", err)
	}

	return filename, nil
}

func (db *Database) InsertMachine(ctx context.Context, machine *domain.Machine) (int64, error) {
	id, err := db.Q.InsertMachine(ctx, sqlc.InsertMachineParams{
		Location: machine.Location,
		Rows:     int32(machine.Rows),
		Columns:  int32(machine.Cols),
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) InsertBookMetric(ctx context.Context, bookID, machineID int64, qr bool, source, sessionID string) (int64, error) {
	sourceText := pgtype.Text{String: source, Valid: source != ""}
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		sid = "unknown"
	}

	metricID, err := db.Q.InsertBookMetric(ctx, sqlc.InsertBookMetricParams{
		BookID:    bookID,
		MachineID: machineID,
		Date: pgtype.Date{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		Qr:        qr,
		Source:    sourceText,
		SessionID: sid,
	})
	if err != nil {
		return 0, err
	}

	return metricID, nil
}

func (db *Database) InsertMachineMetric(ctx context.Context, machineID int64, qr bool, source string, admin bool, sessionID string) (int64, error) {
	s := strings.TrimSpace(strings.ToLower(source))
	if s == "" {
		s = "unknown"
	}
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		sid = "unknown"
	}

	metricID, err := db.Q.InsertMachineMetric(ctx, sqlc.InsertMachineMetricParams{
		MachineID: machineID,
		Qr:        qr,
		Source:    s,
		Admin:     admin,
		SessionID: sid,
	})
	if err != nil {
		return 0, err
	}

	return metricID, nil
}

func (db *Database) LoadMachine(ctx context.Context, machineID int64, books []domain.BookMachine) error {
	tx, err := db.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to open database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Q.WithTx(tx)

	if err := qtx.DeleteBookMachineByMachineID(ctx, machineID); err != nil {
		return fmt.Errorf("failed to clear existing machine load rows: %w", err)
	}

	for _, book := range books {
		err := qtx.InsertBookMachine(ctx, sqlc.InsertBookMachineParams{
			MachineID: machineID,
			BookID:    book.BookID,
			Row:       int32(book.Row),
			Col:       int32(book.Col),
		})
		if err != nil {
			return fmt.Errorf("failed to insert book_machine row: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
