package repository

import "fmt"

func (conn *Database) CreateTables() error {
	_, err := conn.Conn.NewCreateTable().Model((*Book)(nil)).IfNotExists().Exec(conn.Ctx)

	if err != nil {
		return fmt.Errorf("Failed to create table: %w", err)
	}

	_, err = conn.Conn.NewCreateTable().Model((*Machine)(nil)).IfNotExists().Exec(conn.Ctx)

	if err != nil {
		return fmt.Errorf("Failed to create table: %w", err)
	}

	_, err = conn.Conn.NewCreateTable().Model((*BookMachine)(nil)).IfNotExists().Exec(conn.Ctx)

	if err != nil {
		return fmt.Errorf("Failed to create table: %w", err)
	}

	return nil
}
