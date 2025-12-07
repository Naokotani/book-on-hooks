-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE book_machine (
    machine_id  BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
    book_id     BIGINT NOT NULL REFERENCES book(id) ON DELETE CASCADE,

    row         INTEGER NOT NULL,
    col         INTEGER NOT NULL,

    PRIMARY KEY(machine_id, book_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE book_machine;
-- +goose StatementEnd
