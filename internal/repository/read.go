package repository

func (db *Database) GetBookByID(id int64) (*Book, error) {
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

func (db *Database) GetBooks() ([]Book, error) {
	booksSelect := make([]map[string]interface{}, 0)
	err := db.Conn.NewSelect().Model(&Book{}).Limit(100).Scan(db.Ctx, &booksSelect)

	if err != nil {
		return nil, err
	}

	var books []Book

	for _, b := range booksSelect {
		id, _ := b["id"].(int64)
		title, _ := b["title"].(string)
		author, _ := b["author"].(string)
		summary, _ := b["summary"].(string)
		image, _ := b["image"].(string)

		book := Book{
			ID:      id,
			Title:   title,
			Author:  author,
			Summary: summary,
			Image:   image,
		}
		books = append(books, book)
	}

	return books, nil
}

func (db *Database) GetMachines() ([]map[string]interface{}, error) {
	machines := make([]map[string]interface{}, 0)
	var machine Machine
	err := db.Conn.NewSelect().Model(&machine).Limit(100).Scan(db.Ctx, &machines)

	if err != nil {
		return nil, err
	}

	return machines, nil
}
