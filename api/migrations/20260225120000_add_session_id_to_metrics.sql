-- +goose Up
-- +goose StatementBegin
ALTER TABLE book_metrics
ADD COLUMN IF NOT EXISTS session_id TEXT NOT NULL DEFAULT 'unknown';

ALTER TABLE machine_metrics
ADD COLUMN IF NOT EXISTS session_id TEXT NOT NULL DEFAULT 'unknown';

CREATE INDEX IF NOT EXISTS idx_book_metrics_session_id ON book_metrics(session_id);
CREATE INDEX IF NOT EXISTS idx_machine_metrics_session_id ON machine_metrics(session_id);

CREATE INDEX IF NOT EXISTS idx_book_metrics_book_date_session
ON book_metrics(book_id, date, session_id);

CREATE INDEX IF NOT EXISTS idx_machine_metrics_machine_date_session
ON machine_metrics(machine_id, date, session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_machine_metrics_machine_date_session;
DROP INDEX IF EXISTS idx_book_metrics_book_date_session;
DROP INDEX IF EXISTS idx_machine_metrics_session_id;
DROP INDEX IF EXISTS idx_book_metrics_session_id;

ALTER TABLE machine_metrics
DROP COLUMN IF EXISTS session_id;

ALTER TABLE book_metrics
DROP COLUMN IF EXISTS session_id;
-- +goose StatementEnd
