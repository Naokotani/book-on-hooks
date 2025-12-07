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

