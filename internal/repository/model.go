package repository

import (
	"github.com/uptrace/bun"
)

type Book struct {
	bun.BaseModel `bun:"table:book,alias:b"`

	ID     int64  `bun:",pk,autoincrement"`
	Title  string `bun:",notnull"`
	Cover  string `bun:",notnull"`
	Author string `bun:",notnull"`
	Price  string `bun:",notnull"`
}
