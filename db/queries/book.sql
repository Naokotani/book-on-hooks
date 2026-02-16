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
