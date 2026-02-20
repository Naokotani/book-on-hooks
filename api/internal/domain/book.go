package domain

// Domain Models
type Book struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
	Image   string `json:"image"`
	Price   string `json:"price"`
}

type LoadedBook struct {
	Book
	Row int `json:"row"`
	Col int `json:"col"`
}

// Request DTOs
type BookUpdateRequest struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
	Price   string `json:"price"`
}

// Response DTOs
type BookLocation struct {
	BookID    int64             `json:"book_id"`
	Title     string            `json:"title"`
	Author    string            `json:"author"`
	Summary   string            `json:"summary"`
	Image     string            `json:"image"`
	Price     string            `json:"price"`
	Locations []MachineLocation `json:"locations"`
}
