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

-- name: GetFinancialLog :one
SELECT * FROM financial_logs
WHERE id = ? LIMIT 1;

-- name: ListFinancialLogs :many
SELECT * FROM financial_logs
ORDER BY created_at DESC;

-- name: CreateFinancialLog :one
INSERT INTO financial_logs (
  investiment, amount, currency, created_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateFinancialLog :one
UPDATE financial_logs
set
    investiment = ?,
    amount = ?,
    currency = ?,
    created_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteFinancialLog :one
DELETE FROM financial_logs
WHERE id = ?
RETURNING *;
