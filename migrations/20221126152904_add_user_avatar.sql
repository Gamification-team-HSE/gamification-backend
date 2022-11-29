-- +goose Up
-- +goose StatementBegin
alter table users
    add column avatar text,
    add column "name" text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    drop column avatar,
    drop column "name";
-- +goose StatementEnd
