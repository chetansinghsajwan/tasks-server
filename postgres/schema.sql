CREATE TABLE users (
    id            TEXT PRIMARY KEY,
    full_name     TEXT NOT NULL,
    display_name  TEXT,
    email         TEXT UNIQUE NOT NULL
);

CREATE TABLE lists (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id    TEXT NOT NULL REFERENCES users(id),
    name        TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tasks (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    list_id      BIGINT NOT NULL REFERENCES lists(id),
    title        TEXT NOT NULL,
    description  TEXT,
    priority     INTEGER,
    due_date     TIMESTAMP,
    assignee     TEXT REFERENCES users(id),
    labels       TEXT[]
);

CREATE TABLE list_permissions (
    user_id            TEXT REFERENCES users(id),
    list_id            BIGINT REFERENCES lists(id),
    can_read_tasks     BOOLEAN DEFAULT FALSE,
    can_update_tasks   BOOLEAN DEFAULT FALSE,
    can_create_tasks   BOOLEAN DEFAULT FALSE,
    can_delete_tasks   BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id, list_id)
);
