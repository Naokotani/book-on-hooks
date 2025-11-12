package repository

import "errors"

func (conn *Database) CreateTables() error {
	_, err := conn.Conn.NewCreateTable().Model((*Book)(nil)).IfNotExists().Exec(conn.Ctx)

	if err != nil {
		return errors.New("Failed to create table.\n%s\n.")
	}

	return nil
}
