-- name: CreateNewDeployment :one
INSERT INTO deployments(git_url)
values ($1)
Returning *;