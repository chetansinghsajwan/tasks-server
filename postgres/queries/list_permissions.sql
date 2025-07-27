-- name: GetListPermission :one
SELECT * FROM list_permissions
WHERE user_id = @user_id AND list_id = @list_id;

-- name: GetListPermissions :many
SELECT * FROM list_permissions
WHERE (user_id, list_id) IN (
    SELECT unnest(@user_ids::text[]), unnest(@list_ids::bigint[])
);

-- name: GetListPermissionWhere :one
SELECT * FROM list_permissions
WHERE @filter
LIMIT 1;

-- name: GetListPermissionsWhere :many
SELECT * FROM list_permissions
WHERE @filter
LIMIT @count;

-- name: CreateListPermission :exec
INSERT INTO list_permissions (
    user_id,
    list_id,
    can_read_tasks,
    can_update_tasks,
    can_create_tasks,
    can_delete_tasks
) VALUES (
    @user_id,
    @list_id,
    @can_read_tasks,
    @can_update_tasks,
    @can_create_tasks,
    @can_delete_tasks
);

-- name: UpdateListPermission :exec
UPDATE list_permissions SET
    can_read_tasks   = @can_read_tasks,
    can_update_tasks = @can_update_tasks,
    can_create_tasks = @can_create_tasks,
    can_delete_tasks = @can_delete_tasks
WHERE user_id = @user_id AND list_id = @list_id;

-- name: UpdateListPermissions :exec
UPDATE list_permissions SET
    can_read_tasks   = @can_read_tasks,
    can_update_tasks = @can_update_tasks,
    can_create_tasks = @can_create_tasks,
    can_delete_tasks = @can_delete_tasks
WHERE (user_id, list_id) IN (
    SELECT unnest(@user_ids::text[]), unnest(@list_ids::bigint[])
);

-- name: DeleteListPermission :exec
DELETE FROM list_permissions
WHERE user_id = @user_id AND list_id = @list_id;

-- name: DeleteListPermissions :exec
DELETE FROM list_permissions
WHERE (user_id, list_id) IN (
    SELECT unnest(@user_ids::text[]), unnest(@list_ids::bigint[])
);
