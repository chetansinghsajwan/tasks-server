-- name: GetUser :one
SELECT * FROM users
WHERE id = @id;

-- name: GetUsers :many
SELECT * FROM users
WHERE id = ANY(@ids::text[]);

-- name: GetUserWhere :one
SELECT * FROM users
WHERE @filter
LIMIT 1;

-- name: GetUsersWhere :many
SELECT * FROM users
WHERE @filter
LIMIT @count;

-- name: CreateUser :exec
INSERT INTO users (
    id,
    full_name,
    display_name,
    email
) VALUES (
    @id,
    @full_name,
    @display_name,
    @email
);

-- name: UpdateUser :exec
UPDATE users SET
    full_name = @full_name,
    display_name = @display_name,
    email = @email
WHERE id = @id;

-- name: UpdateUsers :exec
UPDATE users SET
    full_name = @full_name,
    display_name = @display_name,
    email = @email
WHERE id = ANY(@ids::text[]);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id;

-- name: DeleteUsers :exec
DELETE FROM users
WHERE id = ANY(@ids::text[]);
