package repository

func (db *Database) GetBookByID(id int) (*Book, error) {
	book := new(Book)
	err := db.Conn.NewSelect().Model(book).Where("id = ?", id).Scan(db.Ctx)

	if err != nil {
		return book, err
	}
	return book, nil
}

func (db *Database) GetBookByRowAndCol(row, col int) (*Book, error) {
	book := new(Book)
	err := db.Conn.NewSelect().Model(book).
		Where("row = ?", row).
		Where("col = ?", col).
		Scan(db.Ctx)

	if err != nil {
		return book, err
	}
	return book, nil
}
