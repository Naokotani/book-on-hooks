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
		Books: []domain.LoadedBook{},
	}

	for _, row := range rows {
		if !row.BookID.Valid {
			continue
		}

		loaded := domain.LoadedBook{
			Book: domain.Book{
				ID:      row.BookID.Int64,
				Title:   row.Title.String,
				Author:  row.Author.String,
				Summary: row.Summary.String,
				Image:   row.Image.String,
				Price:   row.Price.String,
			},
		}
		if row.Row.Valid {
			loaded.Row = int(row.Row.Int32)
		}
		if row.Col.Valid {
			loaded.Col = int(row.Col.Int32)
		}

		out.Books = append(out.Books, loaded)
	}

	return out, nil
}
