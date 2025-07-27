-- name: GetList :one
SELECT * FROM lists
WHERE id = @id;

-- name: GetLists :many
SELECT * FROM lists
WHERE id = ANY(@ids::bigint[]);

-- name: GetListWhere :one
SELECT * FROM lists
WHERE @filter
LIMIT 1;

-- name: GetListsWhere :many
SELECT * FROM lists
WHERE @filter
LIMIT @count;

-- name: CreateList :one
INSERT INTO lists (owner_id, name)
VALUES (@owner_id, @name)
RETURNING id;

-- name: UpdateList :exec
UPDATE lists SET
    owner_id = @owner_id,
    name = @name
WHERE id = @id;

-- name: UpdateLists :exec
UPDATE lists SET
    owner_id = @owner_id,
    name = @name
WHERE id = ANY(@ids::bigint[]);

-- name: DeleteList :exec
DELETE FROM lists
WHERE id = @id;

-- name: DeleteLists :exec
DELETE FROM lists
WHERE id = ANY(@ids::bigint[]);
