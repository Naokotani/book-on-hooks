-- name: InsertMachine :one
INSERT INTO machine (
    location,
    rows,
    columns
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: ListMachines :many
SELECT *
FROM machine
ORDER BY id
LIMIT 100;
