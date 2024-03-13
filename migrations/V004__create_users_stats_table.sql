-- +goose Up
-- +goose StatementBegin
create table users_stats (
    user_id bigint references users ON DELETE CASCADE,
    stat_id bigint references stats ON DELETE CASCADE,
    "value" bigint,
    updated_at timestamp,
    created_at timestamp not null,
    PRIMARY KEY (user_id, stat_id)
);
-- +goose StatementEnd
