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
