-- name: InsertBookMachine :exec
INSERT INTO book_machine (
    machine_id,
    book_id,
    row,
    col
) VALUES ($1, $2, $3, $4);

-- name: UpdateBookPosition :exec
UPDATE book_machine
SET row = $2,
    col = $3
WHERE book_id = $1;

-- name: GetBookByRowAndCol :one
SELECT b.*
FROM book_machine bm
JOIN book b ON b.id = bm.book_id
WHERE bm.row = $1
  AND bm.col = $2
LIMIT 1;

-- name: DeleteBookMachineByMachineID :exec
DELETE FROM book_machine
WHERE machine_id = $1;

-- name: GetMachineWithBooks :many
SELECT
    m.id AS machine_id,
    m.location,
    m.rows AS machine_rows,
    m.columns AS machine_columns,
    b.id AS book_id,
    b.title,
    b.author,
    b.summary,
    b.image,
    b.price,
    bm.row,
    bm.col
FROM machine m
LEFT JOIN book_machine bm ON bm.machine_id = m.id
LEFT JOIN book b ON b.id = bm.book_id
WHERE m.id = $1
ORDER BY bm.row, bm.col;
