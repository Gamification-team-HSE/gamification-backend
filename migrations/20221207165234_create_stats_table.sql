-- +goose Up
-- +goose StatementBegin
create table stats (
    id bigserial primary key,
    "name" text unique not null,
    description text,
    created_at timestamp not null,
    start_at timestamp not null,
    period text not null,
    seq_period text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stats;
-- +goose StatementEnd
