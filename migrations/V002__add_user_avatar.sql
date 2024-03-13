-- +goose Up
-- +goose StatementBegin
alter table users
    add column avatar text,
    add column "name" text;
-- +goose StatementEnd
