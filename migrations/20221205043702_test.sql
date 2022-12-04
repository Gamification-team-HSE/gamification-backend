-- +goose Up
-- +goose StatementBegin
create table test (
    id bigserial
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table test;
-- +goose StatementEnd
