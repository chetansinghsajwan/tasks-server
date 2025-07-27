-- name: GetTask :one
SELECT * FROM tasks
WHERE id = @id;

-- name: GetTasks :many
SELECT * FROM tasks
WHERE id = ANY(@ids::int[]);

-- name: GetTaskWhere :one
SELECT * FROM tasks
WHERE @filter
LIMIT 1;

-- name: GetTasksWhere :many
SELECT * FROM tasks
WHERE @filter
LIMIT @count;

-- name: CreateTask :one
INSERT INTO tasks (title, description, priority, due_date, assignee, labels, list_id)
VALUES (@title, @description, @priority, @due_date, @assignee, @labels, @list_id)
RETURNING id;

-- name: UpdateTask :exec
UPDATE tasks SET
    title = @title,
    description = @description,
    priority = @priority,
    due_date = @due_date,
    assignee = @assignee,
    labels = @labels,
    list_id = @list_id
WHERE id = @id;

-- name: UpdateTasks :exec
UPDATE tasks SET
    title = @title,
    description = @description,
    priority = @priority,
    due_date = @due_date,
    assignee = @assignee,
    labels = @labels,
    list_id = @list_id
WHERE id = ANY(@ids::int[]);

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = @id;

-- name: DeleteTasks :exec
DELETE FROM tasks
WHERE id = ANY(@ids::int[]);
