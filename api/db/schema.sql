-- BOOK TABLE
CREATE TABLE book (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    author      TEXT NOT NULL,
    summary     TEXT NOT NULL,
    image       TEXT NOT NULL,
    price       TEXT NOT NULL
);

-- MACHINE TABLE
CREATE TABLE machine (
    id          BIGSERIAL PRIMARY KEY,
    location    TEXT NOT NULL,
    rows        INTEGER NOT NULL,
    columns     INTEGER NOT NULL
);

-- BOOK_MACHINE (join table)
CREATE TABLE book_machine (
    machine_id  BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
    book_id     BIGINT NOT NULL REFERENCES book(id) ON DELETE CASCADE,

    row         INTEGER NOT NULL,
    col         INTEGER NOT NULL,

    PRIMARY KEY(machine_id, book_id)
);

-- BOOK_METRICS TABLE
CREATE TABLE book_metrics (
    id          BIGSERIAL PRIMARY KEY,
    book_id     BIGINT NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    machine_id  BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
    date        DATE NOT NULL DEFAULT CURRENT_DATE,
    qr          BOOLEAN NOT NULL DEFAULT FALSE,
    source      TEXT,
    session_id  TEXT NOT NULL DEFAULT 'unknown'
);

-- MACHINE_METRICS TABLE
CREATE TABLE machine_metrics (
    id          BIGSERIAL PRIMARY KEY,
    machine_id  BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
    date        DATE NOT NULL DEFAULT CURRENT_DATE,
    qr          BOOLEAN NOT NULL DEFAULT FALSE,
    source      TEXT NOT NULL DEFAULT 'unknown',
    admin       BOOLEAN NOT NULL DEFAULT FALSE,
    session_id  TEXT NOT NULL DEFAULT 'unknown'
);
