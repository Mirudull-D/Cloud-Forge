-- name: CreateNewDeployment :one
INSERT INTO deployments(git_url)
values ($1)
Returning *;

-- name: GetAllDeployments :many
SELECT *
FROM deployments ;

-- name: GetDeploymentById :one
SELECT *
FROM deployments
WHERE id = $1;

-- name: ChangeDeploymentStatus :one
UPDATE deployments
SET status = 'building'
WHERE id = $1
Returning *;