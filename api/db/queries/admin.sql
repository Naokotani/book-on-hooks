-- name: GetTableCounts :one
SELECT
    (SELECT COUNT(*) FROM book)::bigint AS books_count,
    (SELECT COUNT(*) FROM machine)::bigint AS machines_count,
    (SELECT COUNT(*) FROM book_machine)::bigint AS book_machine_count;

-- name: TruncateCoreTables :exec
TRUNCATE TABLE
    book_machine,
    book,
    machine
RESTART IDENTITY CASCADE;
