-- name: GetCronJob :one
SELECT * FROM crons
WHERE id = ? LIMIT 1;

-- name: ListCronJobs :many
SELECT * FROM crons
ORDER BY name;

-- name: CreateCronJob :one
INSERT INTO crons (
  name, expression, script
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateCronJob :exec
UPDATE crons
set
    name = ?,
    expression = ?,
    script = ?
WHERE id = ?;

-- name: DeleteCronJob :exec
DELETE FROM crons
WHERE id = ?
RETURNING *;
