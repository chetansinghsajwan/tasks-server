create or replace function reinitialize_schema()
returns void as $$
begin
    drop table if exists list_access;
    drop table if exists tasks;
    drop table if exists lists;
    drop table if exists secrets;
    drop table if exists users;
    drop type if exists secret_scopes;
    drop type if exists list_access_type;
    drop type if exists user_id_type;
    drop type if exists list_id_type;
    drop type if exists list_name_type;

    create domain user_id_type as varchar(30);
    create domain list_id_type as varchar(30);
    create domain list_name_type as varchar(60);

    create table users (
        id            user_id_type primary key,
        email         text unique not null,
        full_name     text not null,
        display_name  text,

        constraint users_id_validation check (
            length(trim(id)) != 0 and
            id ~ '^[a-z0-9][a-z0-9-]{0,28}[a-z0-9]$' and
            id !~ '--+'
        ),

        constraint users_email_validation check (
            length(trim(email)) != 0 and
            email ~* '^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$'
        ),

        constraint users_full_name_validation check (
            length(trim(full_name)) != 0
        ),

        constraint users_display_name_validation check (
            display_name is null or length(trim(display_name)) != 0
        )
    );

    create type secret_scopes as enum (
        'user-login'
    );

    create table secrets (
        key           text not null,
        scope         secret_scopes not null,
        value         text not null,

        primary key (key, scope),

        constraint secrets_key_validation check (
            length(trim(key)) != 0
        ),

        constraint secrets_value_validation check (
            length(trim(value)) != 0
        )
    );

    create table lists (
        id           bigint generated always as identity primary key,
        name         list_name_type unique not null,
        owner_id     user_id_type not null references users(id),

        constraint lists_id_validation check (
            length(trim(id)) != 0
        )
    );

    create table tasks (
        id           bigint generated always as identity primary key,
        list_id      list_id_type not null references lists(id),
        title        text not null,
        description  text,
        priority     integer,
        due_date     timestamp,
        assignee     user_id_type references users(id),
        labels       text[],

        constraint tasks_title_validation check (
            length(trim(title)) != 0
        ),

        constraint tasks_priority_validation check (
            priority is null or priority > 0
        )
    );

    create type list_access_type as enum (
	    "owner",
	    "add-task",
	    "read-task",
	    "write-task",
	    "remove-task",
	    "add-access",
	    "read-access",
	    "remove-access"
    );

    create table list_access (
        user_id            user_id_type not null references users(id),
        list_id            list_id_type not null references lists(id),
        access             list_access_type not null,

        primary key (user_id, list_id, access)
    );
end;
$$ language plpgsql;

select reinitialize_schema();
