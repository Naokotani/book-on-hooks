package main

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID    int64  `bun:",pk,autoincrement"`
	Name  string `bun:",notnull"`
	Email string `bun:",unique"`
}
