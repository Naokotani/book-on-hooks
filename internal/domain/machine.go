package domain

type Machine struct {
	ID       int64
	Location string
	Rows     int
	Columns  int
}

type BookMachine struct {
	MachineID int64
	BookID    int64

	Row int
	Col int
}

type MachineWithBooks struct {
	Machine Machine
	Books   []Book
}
