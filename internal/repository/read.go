package repository

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/sqlc"
	"context"
)

func (db *Database) GetBookByID(ctx context.Context, id int64) (*domain.Book, error) {
	book, err := db.Q.GetBookByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &domain.Book{
		ID:      book.ID,
		Title:   book.Title,
		Author:  book.Author,
		Summary: book.Summary,
		Image:   book.Image,
		Price:   book.Price,
	}, nil
}

func (db *Database) GetBookByRowAndCol(ctx context.Context, row, col int) (*domain.Book, error) {
	book, err := db.Q.GetBookByRowAndCol(ctx,
		sqlc.GetBookByRowAndColParams{
			Row: int32(row),
			Col: int32(col),
		})

	if err != nil {
		return nil, err
	}
	return mapBook(book), nil
}

func (db *Database) GetBooks(ctx context.Context) ([]domain.Book, error) {
	books, err := db.Q.ListBooks(ctx)

	if err != nil {
		return nil, err
	}

	return mapBooks(books), nil
}

func (db *Database) GetMachines(ctx context.Context) ([]domain.Machine, error) {
	machines, err := db.Q.ListMachines(ctx)
	if err != nil {
		return nil, err
	}

	return mapMachines(machines), nil
}

func (db *Database) GetMachineById(ctx context.Context, id int64) (*domain.Machine, error) {
	machine, err := db.Q.GetMachineById(ctx, id)
	if err != nil {
		return &domain.Machine{}, err
	}

	return mapMachine(machine), nil
}
