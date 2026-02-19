-- +goose Up
-- +goose StatementBegin
CREATE TABLE book_metrics (
    id          BIGSERIAL PRIMARY KEY,
    book_id     BIGINT NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    machine_id  BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
    date        DATE NOT NULL DEFAULT CURRENT_DATE,
    qr          BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE machine_metrics (
    id          BIGSERIAL PRIMARY KEY,
    machine_id  BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
    date        DATE NOT NULL DEFAULT CURRENT_DATE,
    qr          BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS machine_metrics;
DROP TABLE IF EXISTS book_metrics;
-- +goose StatementEnd
