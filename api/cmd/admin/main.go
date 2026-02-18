package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"booksonhooks.ca/internal/repository"
)

const usage = `Usage:
  go run ./cmd/admin counts
  go run ./cmd/admin reset-database`

func main() {
	if len(os.Args) < 2 {
		exitUsage()
	}

	db, err := repository.GetDatabaseConnection()
	if err != nil {
		fatalf("failed to create database connection: %v\n", err)
	}
	defer db.Db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch os.Args[1] {
	case "counts":
		countsCommand(ctx, db)

	case "reset-database":
		resetDatabaseCommand(ctx, db)

	default:
		exitUsage()
	}
}

func countsCommand(ctx context.Context, db *repository.Database) {
	if len(os.Args) != 2 {
		exitUsage()
	}

	row, err := db.Q.GetTableCounts(ctx)
	if err != nil {
		fatalf("failed to get table counts: %v\n", err)
	}

	fmt.Printf("books=%d machines=%d book_machine=%d\n",
		row.BooksCount, row.MachinesCount, row.BookMachineCount)
}

func resetDatabaseCommand(ctx context.Context, db *repository.Database) {
	if err := db.Q.TruncateCoreTables(ctx); err != nil {
		fatalf("failed to truncate core tables: %v\n", err)
	}
	row, err := db.Q.GetTableCounts(ctx)
	if err != nil {
		fatalf("database reset completed, but failed to read table counts: %v\n", err)
	}

	fmt.Printf("database reset complete: books=%d machines=%d book_machine=%d\n",
		row.BooksCount, row.MachinesCount, row.BookMachineCount)
}

func exitUsage() {
	fmt.Fprintln(os.Stderr, usage)
	os.Exit(2)
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
