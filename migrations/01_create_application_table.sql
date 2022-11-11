-- +migrate Up
create table if not exists application
(
    id              bigserial    not null primary key,
    requestId       bigserial,
    checkedId       bigserial,
    person          varchar(32),
    customerType    varchar(32),
    CustomerName    varchar(32),
    FilePath        varchar(32),
    CourtName       varchar(32),
    JudgeName       varchar(32),
    DecisionNumber  varchar(32),
    DecisionDate    varchar(32),
    IsChecked       boolean,
    Note            varchar(32),
    Status          varchar(32),
    Deadline        varchar(32),
    AssigneeId      bigserial,
    Priority        varchar(32),
    MailSent        boolean,
    AssigneeName    varchar(32),
    BeginDate       varchar(32),
    EndDate         varchar(32),
    CreatedAt       timestamp    not null default now()

);
