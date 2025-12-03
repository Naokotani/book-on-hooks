package repository

import (
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

func (db *Database) InsertBook(book *Book, file multipart.File, header *multipart.FileHeader) (int64, error) {

	tx, err := db.Conn.BeginTx(db.Ctx, &sql.TxOptions{})

	_, err = tx.NewInsert().Model(book).Exec(db.Ctx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	filename := strconv.FormatInt(book.ID, 10) + "_" + header.Filename

	dst, err := os.Create(filepath.Join("./images/covers", filename))
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.NewUpdate().
		Model(&Book{}).
		Set("cover = ?", filename).
		Where("id = ?", book.ID).
		Exec(db.Ctx)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return book.ID, nil
}

func (db *Database) InsertMachine(machine *Machine) error {
	_, err := db.Conn.NewInsert().Model(machine).Exec(db.Ctx)
	if err != nil {
		return err
	}

	return nil
}
