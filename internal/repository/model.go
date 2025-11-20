package repository

import (
	"github.com/uptrace/bun"
)

type Book struct {
	bun.BaseModel `bun:"table:book,alias:b"`

	ID           int64          `bun:",pk,autoincrement"`
	Title        string         `bun:",notnull"`
	Author       string         `bun:",notnull"`
	Summary      string         `bun:",notnull"`
	Image        string         `bun:",notnull"`
	Price        float64        `bun:",notnull"`
	BookMachines []*BookMachine `bun:"rel:m2m:book_machine,join:Book=Machine"`
}

type Machine struct {
	bun.BaseModel `bun:"table:machine,alias:m"`

	ID           int64          `bun:",pk,autoincrement"`
	Provice      string         `bun:",notnull"`
	City         string         `bun:",notnull"`
	Street       string         `bun:",notnull"`
	PostalCode   string         `bun:",notnull"`
	BookMachines []*BookMachine `bun:"rel:m2m:book_machine,join:Machine=Book"`
}

type BookMachine struct {
	MachineID int64 `bun:",pk"`
	BookID    int64 `bun:",pk"`

	Row int
	Col int

	Machine *Machine `bun:"rel:belongs-to,join:machine_id=id"`
	Book    *Book    `bun:"rel:belongs-to,join:book_id=id"`
}
