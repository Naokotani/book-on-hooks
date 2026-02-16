package domain

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
	Machine Machine `json:"machine"`
	Books   []Book  `json:"books"`
}
