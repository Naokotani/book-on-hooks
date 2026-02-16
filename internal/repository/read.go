package repository

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/sqlc"
	"context"

	"github.com/jackc/pgx/v5"
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

func (db *Database) GetMachineWithBooks(ctx context.Context, id int64) (*domain.MachineWithBooks, error) {
	rows, err := db.Q.GetMachineWithBooks(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, pgx.ErrNoRows
	}

	out := &domain.MachineWithBooks{
		Machine: domain.Machine{
			ID:       rows[0].MachineID,
			Location: rows[0].Location,
			Rows:     int(rows[0].MachineRows),
			Columns:  int(rows[0].MachineColumns),
		},
		Books: []domain.Book{},
	}

	for i := 0; i < len(rows) && len(rows) != 0; i++ {
		if rows[i].BookID.Valid {
			book := domain.Book{
				ID:      rows[i].BookID.Int64,
				Title:   rows[i].Title.String,
				Author:  rows[i].Author.String,
				Summary: rows[i].Summary.String,
				Image:   rows[i].Image.String,
				Price:   rows[i].Price.String,
			}
			out.Books = append(out.Books, book)
		}
	}

	return out, nil
}
