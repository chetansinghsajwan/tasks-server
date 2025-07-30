-- name: CreateSecret :exec
INSERT INTO secrets (key, scope, value)
VALUES (@key, @scope, @pass);

-- name: GetSercert :one
SELECT value
FROM secrets
WHERE key = @key AND scope = @scope;

-- name: UpdateSecret :exec
UPDATE secrets
SET value = @value
WHERE key = @key AND scope = @scope;

-- name: DeleteSecret :exec
UPDATE secrets
SET value = @value
WHERE key = @key AND scope = @scope;
