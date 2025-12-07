package repository

import (
	"booksonhooks.ca/internal/sqlc"
	"context"
)

func (db *Database) GetBookByID(ctx context.Context, id int64) (*Book, error) {
	q := sqlc.New(db.Db)

	book, err := q.GetBookByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &Book{
		book.ID,
		book.Title,
		book.Author,
		book.Summary,
		book.Image,
		book.Price,
	}, nil
}

func (db *Database) GetBookByRowAndCol(ctx context.Context, row, col int) (*Book, error) {
	q := sqlc.New(db.Db)

	book, err := q.GetBookByRowAndCol(ctx,
		sqlc.GetBookByRowAndColParams{
			Row: int32(row),
			Col: int32(col),
		})

	if err != nil {
		return nil, err
	}
	return mapBook(book), nil
}

func (db *Database) GetBooks(ctx context.Context) ([]Book, error) {
	q := sqlc.New(db.Db)

	books, err := q.ListBooks(ctx)

	if err != nil {
		return nil, err
	}

	return mapBooks(books), nil
}

func (db *Database) GetMachines(ctx context.Context) ([]Machine, error) {
	q := sqlc.New(db.Db)

	machines, err := q.ListMachines(ctx)
	if err != nil {
		return nil, err
	}

	return mapMachines(machines), nil
}
