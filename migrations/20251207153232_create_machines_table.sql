-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE machine (
    id          BIGSERIAL PRIMARY KEY,
    location    TEXT NOT NULL,
    rows        INTEGER NOT NULL,
    columns     INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE machine;
-- +goose StatementEnd
