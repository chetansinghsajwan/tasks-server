CREATE OR REPLACE FUNCTION reinitialize_schema()
RETURNS void AS $$
BEGIN
    DROP TABLE IF EXISTS list_access;
    DROP TABLE IF EXISTS tasks;
    DROP TABLE IF EXISTS lists;
    DROP TABLE IF EXISTS secrets;
    DROP TABLE IF EXISTS users;
    DROP TYPE IF EXISTS secret_scopes;

    CREATE TABLE users (
        id            TEXT PRIMARY KEY,
        full_name     TEXT NOT NULL,
        display_name  TEXT,
        email         TEXT UNIQUE NOT NULL
    );

    CREATE TYPE secret_scopes AS ENUM (
        'user-login'
    );

    CREATE TABLE secrets (
        key           TEXT NOT NULL,
        scope         secret_scopes NOT NULL,
        value         TEXT NOT NULL,
        PRIMARY KEY (key, scope)
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

    CREATE TABLE list_access (
        user_id            TEXT REFERENCES users(id),
        list_id            BIGINT REFERENCES lists(id),
        can_read_tasks     BOOLEAN NOT NULL DEFAULT FALSE,
        can_update_tasks   BOOLEAN NOT NULL DEFAULT FALSE,
        can_create_tasks   BOOLEAN NOT NULL DEFAULT FALSE,
        can_delete_tasks   BOOLEAN NOT NULL DEFAULT FALSE,
        PRIMARY KEY (user_id, list_id)
    );
END;
$$ LANGUAGE plpgsql;

SELECT reinitialize_schema();
