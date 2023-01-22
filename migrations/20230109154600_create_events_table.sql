-- +goose Up
-- +goose StatementBegin
create table event(
    id          bigserial primary key,
    "name"      text unique not null,
    description text,
    "image"     text,
    created_at  timestamp not null,
    start_at    timestamp not null,
    end_at      timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table event;
-- +goose StatementEnd
