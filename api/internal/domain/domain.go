package domain

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

type Machine struct {
	ID       int64  `json:"id"`
	Location string `json:"location"`
	Rows     int    `json:"rows"`
	Columns  int    `json:"columns"`
}

type BookMachine struct {
	MachineID int64 `json:"machine_id"`
	BookID    int64 `json:"book_id"`

	Row int `json:"row"`
	Col int `json:"col"`
}

type MachineWithBooks struct {
	Machine Machine      `json:"machine"`
	Books   []LoadedBook `json:"books"`
}

type MachineLocation struct {
	MachineID int64  `json:"machine_id"`
	Location  string `json:"location"`
}

type BookLocation struct {
	BookID    int64             `json:"book_id"`
	Title     string            `json:"title"`
	Author    string            `json:"author"`
	Summary   string            `json:"summary"`
	Image     string            `json:"image"`
	Price     string            `json:"price"`
	Locations []MachineLocation `json:"locations"`
}
