create or replace function reinitialize_schema()
returns void as $$
begin
    drop table if exists list_access cascade;
    drop table if exists tasks cascade;
    drop table if exists lists cascade;
    drop table if exists users cascade;
    drop table if exists user_secrets cascade;
    drop type if exists list_access_type;
    drop type if exists user_id_type;
    drop type if exists list_name_type;

    create domain user_id_type as varchar(30);
    create domain list_name_type as varchar(60);

    create table users (
        id            user_id_type primary key,
        email         text unique not null,
        full_name     text not null,
        display_name  text,

        constraint users_id_validation check (

            -- ensure no empty whitespace value
            length(trim(id)) != 0

            -- ensure lowercase alphanumeric and hyphen only and between 2 to 30 in length
            and id ~ '^[a-z0-9][a-z0-9-]{0,28}[a-z0-9]$'

            -- ensure no multiple hyphens together
            and id !~ '--+'
        ),

        constraint users_email_validation check (

            -- ensure no empty whitespace value
            length(trim(email)) != 0

            -- ensure valid email format:
            -- - lowercase only
            -- - local part starts/ends with alnum, allows ._%+-
            -- - domain part allows alnum, dot, hyphen
            -- - TLD at least 2 chars
            and email ~ '^[a-z0-9](?:[a-z0-9._%+-]*[a-z0-9])?@[a-z0-9](?:[a-z0-9.-]*[a-z0-9])?\.[a-z]{2,}$'
        ),

        constraint users_full_name_validation check (

            -- ensure no empty whitespace value
            length(trim(full_name)) != 0

            -- ensure no leading and trailing spaces
            and full_name = btrim(full_name)

            -- forbid tabs, newlines, carriage returns
            and full_name !~ '[\t\n\r]'
        ),

        constraint users_display_name_validation check (

            display_name is null or (

                -- ensure no empty whitespace value
                length(trim(display_name)) != 0

                -- ensure no leading and trailing spaces
                and display_name = btrim(display_name)

                -- forbid tabs, newlines, carriage returns
                and full_name !~ '[\t\n\r]'
            )
        )
    );

    create table user_secrets (
        id          user_id_type not null,
        value       text not null,

        constraint user_secrets_pk primary key (id),

        constraint user_secrets_id_fk foreign key (id)
            references users(id) on delete cascade,

        constraint user_secrets_value_format_check check (
            length(trim(value)) != 0
        )
    );

    create table lists (
        id           bigint generated always as identity primary key,
        name         list_name_type unique not null,

        constraint lists_name_validation check (

            -- ensure no empty whitespace characters
            length(trim(name)) != 0
        )
    );

    create table tasks (
        id           bigint generated always as identity primary key,
        list_id      bigint not null references lists(id),
        title        text not null,
        description  text,
        priority     integer,
        due_date     timestamp,
        assignee     user_id_type references users(id),
        labels       text[],

        constraint tasks_title_validation check (

            -- ensure no empty whitespace characters
            length(trim(title)) != 0
        ),

        constraint tasks_priority_validation check (
            priority is null or priority > 0
        )
    );

    create type list_access_type as enum (
	    'owner',
	    'add-task',
	    'read-task',
	    'write-task',
	    'remove-task',
	    'add-access',
	    'read-access',
	    'remove-access'
    );

    create table list_access (
        user_id            user_id_type not null references users(id),
        list_id            bigint not null references lists(id),
        access             list_access_type not null,

        primary key (user_id, list_id, access)
    );
end;
$$ language plpgsql;
