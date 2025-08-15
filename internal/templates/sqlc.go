package templates

// SQLC configuration template
const SqlcConfigTemplate = `version: "2"
sql:
  - engine: "postgresql"
    queries: "./db/queries"
    schema: "./db/migrations"
    gen:
      go:
        package: "sqlc"
        out: "./internal/infrastructure/database/postgres/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_db_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
`

// Sample SQLC query template
const SqlcQueryTemplate = `-- name: Get{{.Entity}} :one
SELECT * FROM {{.Entity | lower}}s WHERE id = $1 LIMIT 1;

-- name: List{{.Entity}}s :many
SELECT * FROM {{.Entity | lower}}s
ORDER BY created_at DESC;

-- name: Create{{.Entity}} :one
INSERT INTO {{.Entity | lower}}s (
  created_at, updated_at
) VALUES (
  $1, $2
)
RETURNING *;

-- name: Update{{.Entity}} :one
UPDATE {{.Entity | lower}}s
SET updated_at = $2
WHERE id = $1
RETURNING *;

-- name: Delete{{.Entity}} :exec
DELETE FROM {{.Entity | lower}}s
WHERE id = $1;
`
