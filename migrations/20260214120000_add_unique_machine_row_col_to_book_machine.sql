-- +goose Up
-- +goose StatementBegin
ALTER TABLE book_machine
ADD CONSTRAINT uq_book_machine_machine_id_row_col UNIQUE (machine_id, row, col);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE book_machine
DROP CONSTRAINT IF EXISTS uq_book_machine_machine_id_row_col;
-- +goose StatementEnd
