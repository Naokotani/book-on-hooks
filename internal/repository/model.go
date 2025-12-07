package repository

type Book struct {
	ID      int64
	Title   string
	Author  string
	Summary string
	Image   string
	Price   string
}

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
