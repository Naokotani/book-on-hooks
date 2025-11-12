package repository

func (db *Database) InsertBook(book *Book) error {

	_, err := db.Conn.NewInsert().Model(book).Exec(db.Ctx)

	if err != nil {
		return err
	}
	return nil
}
