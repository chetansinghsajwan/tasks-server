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

-- name: GetListAccess :one
SELECT * FROM list_access
WHERE user_id = @user_id AND list_id = @list_id;

-- name: SetListAccess :exec
INSERT INTO list_access (user_id, list_id, can_read_tasks, can_update_tasks, can_create_tasks, can_delete_tasks)
VALUES (@user_id, @list_id, @can_read_tasks, @can_update_tasks, @can_create_tasks, @can_delete_tasks)

ON CONFLICT(user_id, list_id) DO

UPDATE SET can_read_tasks = EXCLUDED.can_read_tasks
    , can_update_tasks = EXCLUDED.can_update_tasks
    , can_create_tasks = EXCLUDED.can_create_tasks
    , can_delete_tasks = EXCLUDED.can_delete_tasks;
