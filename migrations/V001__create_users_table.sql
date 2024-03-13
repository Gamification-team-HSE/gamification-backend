-- +goose Up
-- +goose StatementBegin
create table users (
    id bigserial primary key,
    foreign_id text unique,
    email text not null unique,
    created_at timestamp not null,
    deleted_at timestamp,
    "role" text
);
-- +goose StatementEnd
