-- name: InsertBook :one
INSERT INTO book (
    title,
    author,
    summary,
    image,
    price
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id;

-- name: UpdateBookImage :exec
UPDATE book
SET image = $2
WHERE id = $1;

-- name: GetBookByID :one
SELECT *
FROM book
WHERE id = $1;

-- name: ListBooks :many
SELECT *
FROM book
ORDER BY id
LIMIT 100;

-- name: UpdateBook :exec
UPDATE book
SET title = $2,
    author = $3,
    summary = $4,
    price = $5
WHERE id = $1;

-- name: DeleteBook :exec
DELETE FROM book
WHERE id = $1;

-- name: GetBookLocations :many
SELECT
    b.id AS book_id,
    b.title,
    b.author,
    b.summary,
    b.image,
    b.price,
    m.id AS machine_id,
    m.location
FROM book_machine bm
JOIN book b ON b.id = bm.book_id
JOIN machine m ON m.id = bm.machine_id
WHERE bm.book_id = $1
ORDER BY m.location;

-- name: InsertBookMetric :one
INSERT INTO book_metrics (
    book_id,
    machine_id,
    date,
    qr,
    source,
    session_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id;
