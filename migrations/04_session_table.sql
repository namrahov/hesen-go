-- +migrate Up
create table if not exists sessionn
(
    id                 bigserial    not null primary key,
    session_id         varchar(64),
    user_id            varchar(64)
    );
