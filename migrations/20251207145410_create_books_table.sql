-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE book (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    author      TEXT NOT NULL,
    summary     TEXT NOT NULL,
    image       TEXT NOT NULL,
    price       TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE book;
-- +goose StatementEnd
