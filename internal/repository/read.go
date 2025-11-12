package repository

func (db *Database) GetBook(id int) (*Book, error) {
	book := new(Book)
	err := db.Conn.NewSelect().Model(book).Where("id = ?", 1).Scan(db.Ctx)

	if err != nil {
		return book, err
	}
	return book, nil
}
