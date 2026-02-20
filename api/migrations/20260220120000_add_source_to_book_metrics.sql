-- +goose Up
ALTER TABLE book_metrics
ADD COLUMN source TEXT;

-- +goose Down
ALTER TABLE book_metrics
DROP COLUMN IF EXISTS source;
