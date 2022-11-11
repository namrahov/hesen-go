-- +migrate Up
create table if not exists application
(
    id               bigserial    not null primary key,
    request_id        bigserial,
    checked_id        bigserial,
    person           varchar(32),
    customer_type    varchar(32),
    customer_name    varchar(32),
    file_path        varchar(32),
    court_name       varchar(32),
    judge_name       varchar(32),
    decision_number  varchar(32),
    decision_date    timestamp,
    is_checked        boolean,
    note             varchar(32),
    status           varchar(32),
    deadline         timestamp,
    assignee_id       bigserial,
    priority         varchar(32),
    mail_sent         boolean,
    assignee_name     varchar(32),
    begin_date        timestamp,
    end_date          timestamp,
    created_at        timestamp    not null default now()

);
