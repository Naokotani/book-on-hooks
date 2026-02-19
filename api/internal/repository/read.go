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

	return mapBook(book), nil
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

func (db *Database) GetBookLocations(ctx context.Context, bookID int64) (*domain.BookLocation, error) {
	rows, err := db.Q.GetBookLocations(ctx, bookID)
	if err != nil {
		book, err := db.Q.GetBookByID(ctx, bookID)
		if err != nil {
			return nil, err
		}

		return &domain.BookLocation{
			BookID:    book.ID,
			Title:     book.Title,
			Author:    book.Author,
			Summary:   book.Summary,
			Image:     book.Image,
			Price:     book.Price,
			Locations: make([]domain.MachineLocation, 0),
		}, nil
	}

	if len(rows) == 0 {
		return nil, pgx.ErrNoRows
	}

	result := &domain.BookLocation{
		BookID:    rows[0].BookID,
		Title:     rows[0].Title,
		Author:    rows[0].Author,
		Summary:   rows[0].Summary,
		Image:     rows[0].Image,
		Price:     rows[0].Price,
		Locations: make([]domain.MachineLocation, 0, len(rows)),
	}

	for _, row := range rows {
		result.Locations = append(result.Locations, domain.MachineLocation{
			MachineID: row.MachineID,
			Location:  row.Location,
		})
	}

	return result, nil
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

		out.Books = append(out.Books, mapLoadedBook(row))
	}

	return out, nil
}
