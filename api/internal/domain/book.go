package domain

type Book struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
	Image   string `json:"image"`
	Price   string `json:"price"`
}
