-- name: ListTargets :many
SELECT * FROM targets ORDER BY name;

-- name: GetTarget :one
SELECT * FROM targets WHERE id = ?;

-- name: GetTargetByName :one
SELECT * FROM targets WHERE name = ?;

-- name: CreateTarget :one
INSERT INTO targets (name, path) VALUES (?, ?) RETURNING *;

-- name: UpdateTarget :one
UPDATE targets SET name = ?, path = ? WHERE id = ? RETURNING *;

-- name: DeleteTarget :exec
DELETE FROM targets WHERE id = ?;

-- name: CountTargets :one
SELECT COUNT(*) FROM targets;
