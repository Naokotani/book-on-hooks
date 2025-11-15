package repository

func (db *Database) UpdateBookPosition(id, row, col int) error {
	_, err := db.Conn.NewUpdate().
		Model(&Book{}).
		Set("col = ?", col).
		Set("row = ?", row).
		Where("id = ?", id).
		Exec(db.Ctx)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateCoverFilename(id int64, cover string) error {
	_, err := db.Conn.NewUpdate().
		Model(&Book{}).
		Set("cover = ?", cover).
		Where("id = ?", id).
		Exec(db.Ctx)

	if err != nil {
		return err
	}

	return nil
}
