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

-- name: GetMachineById :one
SELECT *
FROM machine
WHERE id = $1;

-- name: UpdateMachine :exec
UPDATE machine
SET location = $2,
    rows = $3,
    columns = $4
WHERE id = $1;

-- name: DeleteMachine :exec
DELETE FROM machine
WHERE id = $1;

-- name: InsertMachineMetric :one
INSERT INTO machine_metrics (
    machine_id,
    date,
    qr,
    source,
    admin,
    session_id
) VALUES (
    $1, NOW()::date, $2, $3, $4, $5
)
RETURNING id;
