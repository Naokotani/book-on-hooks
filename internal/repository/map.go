package repository

import "booksonhooks.ca/internal/sqlc"

func mapBook(book sqlc.Book) *Book {
	return &Book{
		book.ID,
		book.Title,
		book.Author,
		book.Summary,
		book.Image,
		book.Price,
	}
}

func mapBooks(books []sqlc.Book) []Book {
	out := make([]Book, len(books))
	for i, b := range books {
		out[i] = Book{
			b.ID,
			b.Title,
			b.Author,
			b.Summary,
			b.Image,
			b.Price,
		}
	}
	return out
}

func mapMachines(machines []sqlc.Machine) []Machine {
	out := make([]Machine, len(machines))
	for i, m := range machines {
		out[i] = Machine{
			m.ID,
			m.Location,
			int(m.Columns),
			int(m.Rows),
		}
	}
	return out
}
