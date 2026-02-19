package repository

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/sqlc"
)

func mapBook(book sqlc.Book) *domain.Book {
	return &domain.Book{
		ID:      book.ID,
		Title:   book.Title,
		Author:  book.Author,
		Summary: book.Summary,
		Image:   book.Image,
		Price:   book.Price,
	}
}

func mapBooks(books []sqlc.Book) []domain.Book {
	out := make([]domain.Book, len(books))
	for i, b := range books {
		out[i] = domain.Book{
			ID:      b.ID,
			Title:   b.Title,
			Author:  b.Author,
			Summary: b.Summary,
			Image:   b.Image,
			Price:   b.Price,
		}
	}
	return out
}

func mapMachines(machines []sqlc.Machine) []domain.Machine {
	out := make([]domain.Machine, len(machines))
	for i, m := range machines {
		out[i] = domain.Machine{
			ID:       m.ID,
			Location: m.Location,
			Columns:  int(m.Columns),
			Rows:     int(m.Rows),
		}
	}
	return out
}

func mapMachine(machine sqlc.Machine) *domain.Machine {
	return &domain.Machine{
		ID:       machine.ID,
		Location: machine.Location,
		Columns:  int(machine.Columns),
		Rows:     int(machine.Rows),
	}
}

func mapLoadedBook(row sqlc.GetMachineWithBooksRow) domain.LoadedBook {
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

	return loaded
}
