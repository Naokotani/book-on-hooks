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

-- name: GetMetricsTotals :one
SELECT
    (SELECT COUNT(*) FROM book_metrics bm
     WHERE bm.date >= sqlc.arg(start_date)
       AND bm.date < sqlc.arg(end_date)
       AND (sqlc.narg(qr)::boolean IS NULL OR bm.qr = sqlc.narg(qr)::boolean))::bigint AS book_clicks,
    (SELECT COUNT(*) FROM machine_metrics mm
     WHERE mm.date >= sqlc.arg(start_date)
       AND mm.date < sqlc.arg(end_date)
       AND (sqlc.narg(qr)::boolean IS NULL OR mm.qr = sqlc.narg(qr)::boolean))::bigint AS machine_views,
    (SELECT COUNT(DISTINCT bm.session_id) FROM book_metrics bm
     WHERE bm.date >= sqlc.arg(start_date)
       AND bm.date < sqlc.arg(end_date)
       AND (sqlc.narg(qr)::boolean IS NULL OR bm.qr = sqlc.narg(qr)::boolean))::bigint AS unique_sessions;

-- name: GetBookMetricsSummary :many
SELECT
    b.id AS book_id,
    b.title,
    b.author,
    COUNT(bm.id)::bigint AS clicks,
    COUNT(DISTINCT bm.session_id)::bigint AS unique_sessions
FROM book_metrics bm
JOIN book b ON b.id = bm.book_id
WHERE bm.date >= sqlc.arg(start_date)
  AND bm.date < sqlc.arg(end_date)
  AND (sqlc.narg(qr)::boolean IS NULL OR bm.qr = sqlc.narg(qr)::boolean)
GROUP BY b.id, b.title, b.author
ORDER BY clicks DESC, b.title;

-- name: GetMachineMetricsSummary :many
SELECT
    m.id AS machine_id,
    m.location,
    COUNT(DISTINCT mm.id)::bigint AS views,
    COUNT(DISTINCT bm.id)::bigint AS book_clicks,
    COUNT(DISTINCT COALESCE(NULLIF(mm.session_id, 'unknown'), NULLIF(bm.session_id, 'unknown')))::bigint AS unique_sessions
FROM machine m
LEFT JOIN machine_metrics mm ON mm.machine_id = m.id
    AND mm.date >= sqlc.arg(start_date)
    AND mm.date < sqlc.arg(end_date)
    AND (sqlc.narg(qr)::boolean IS NULL OR mm.qr = sqlc.narg(qr)::boolean)
LEFT JOIN book_metrics bm ON bm.machine_id = m.id
    AND bm.date >= sqlc.arg(start_date)
    AND bm.date < sqlc.arg(end_date)
    AND (sqlc.narg(qr)::boolean IS NULL OR bm.qr = sqlc.narg(qr)::boolean)
GROUP BY m.id, m.location
ORDER BY views DESC, book_clicks DESC, m.location;
