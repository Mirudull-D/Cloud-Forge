-- name: CreateNewDeployment :one
INSERT INTO deployments(git_url,status)
values ($1,'queued')
Returning *;

-- name: GetAllDeployments :many
SELECT *
FROM deployments ;

-- name: GetDeploymentById :one
SELECT *
FROM deployments
WHERE id = $1;

-- name: ChangeDeploymentStatusToBuilding :one
UPDATE deployments
SET status = 'building'
WHERE id = $1
Returning *;

-- name: ChangeDeploymentStatusToRunning :one
UPDATE deployments
SET
    status = 'running',
    image_name = $2,
    container_id = $3,
    port = $4,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: ChangeDeploymentStatusToFailed :one
UPDATE deployments
SET
    status = 'failed',
    error_message = $2
WHERE id = $1
    Returning *;

-- name: ChangeDeploymentStatusToStopped :one
UPDATE deployments
SET status = 'stopped'
WHERE id = $1
    Returning *;
-- name: DeleteDeployment :exec
DELETE FROM deployments
WHERE id =$1;

-- name: DeleteContainer :one
UPDATE deployments
SET
    container_id = NULL
WHERE id = $1
    Returning *;