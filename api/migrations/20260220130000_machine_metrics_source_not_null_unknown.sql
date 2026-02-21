-- +goose Up
-- +goose StatementBegin
UPDATE machine_metrics
SET source = 'unknown'
WHERE source IS NULL OR btrim(source) = '';

ALTER TABLE machine_metrics
ALTER COLUMN source SET DEFAULT 'unknown';

ALTER TABLE machine_metrics
ALTER COLUMN source SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE machine_metrics
ALTER COLUMN source DROP NOT NULL;

ALTER TABLE machine_metrics
ALTER COLUMN source DROP DEFAULT;
-- +goose StatementEnd
