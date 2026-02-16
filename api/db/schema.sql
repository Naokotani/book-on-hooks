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
