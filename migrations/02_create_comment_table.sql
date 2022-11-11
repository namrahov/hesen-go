-- +migrate Up
create table if not exists comment
(
    id                 bigserial    not null primary key,
    commentator        varchar(32),
    description        varchar(32),
    comment_type       varchar(32),
    created_at         timestamp    not null default now(),
    updated_at         timestamp    not null default current_timestamp,
    application_id     bigserial
    );

