-- +goose Up
-- +goose StatementBegin
ALTER TABLE machine_metrics
ADD COLUMN IF NOT EXISTS source TEXT;

ALTER TABLE machine_metrics
ADD COLUMN IF NOT EXISTS admin BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_machine_metrics_machine_id ON machine_metrics(machine_id);
CREATE INDEX IF NOT EXISTS idx_machine_metrics_date ON machine_metrics(date);
CREATE INDEX IF NOT EXISTS idx_machine_metrics_source ON machine_metrics(source);
CREATE INDEX IF NOT EXISTS idx_machine_metrics_admin ON machine_metrics(admin);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_machine_metrics_admin;
DROP INDEX IF EXISTS idx_machine_metrics_source;
DROP INDEX IF EXISTS idx_machine_metrics_date;
DROP INDEX IF EXISTS idx_machine_metrics_machine_id;

ALTER TABLE machine_metrics
DROP COLUMN IF EXISTS admin;

ALTER TABLE machine_metrics
DROP COLUMN IF EXISTS source;
-- +goose StatementEnd
