-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists userr
(
    id uuid            DEFAULT uuid_generate_v4() primary key,
    user_name          varchar(32),
    first_name         varchar(32),
    last_name          varchar(32),
    password           bytea        not null,
    created_at         timestamp    not null default now(),
    updated_at         timestamp    not null default current_timestamp
    );
