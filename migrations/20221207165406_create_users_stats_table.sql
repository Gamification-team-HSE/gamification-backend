-- +goose Up
-- +goose StatementBegin
create table users_stats (
    user_id bigint,
    stat_id bigint,
    "value" bigint,
    updated_at bigint,
    created_at timestamp not null,
    PRIMARY KEY (user_id, stat_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users_stats;
-- +goose StatementEnd
