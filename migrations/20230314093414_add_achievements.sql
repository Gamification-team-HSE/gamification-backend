-- +goose Up
-- +goose StatementBegin
create table achievements
(
    id          bigserial primary key,
    "name"      text      not null,
    description text,
    image       text,
    rules       jsonb     not null default '{}',
    end_at      timestamp,
    created_at  timestamp not null default now()
);

create table user_events
(
    user_id    bigint references users on delete cascade,
    event_id   bigint references "event" on delete cascade,
    created_at timestamp not null default now(),
    primary key (user_id, event_id)
);

create table user_achievements
(
    user_id        bigint references users on delete cascade,
    achievement_id bigint references achievements on delete cascade,
    created_at     timestamp not null default now(),
    primary key (user_id, achievement_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_achievements;
drop table user_events;
drop table achievements;
-- +goose StatementEnd
