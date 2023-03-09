package models

import (
	"database/sql"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// DbEvent Select из базы возвращает событие в текущей структуре
type DbEvent struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name" validate:"lte=128"`
	Description sql.NullString `db:"description"`
	Image       sql.NullString `db:"image"`
	CreatedAt   time.Time      `db:"created_at"`
	StartAt     time.Time      `db:"start_at" validate:"required"`
	EndAt       sql.NullTime   `db:"end_at"`
}
type EventTime struct {
	CreatedAt time.Time    `db:"created_at"`
	StartAt   time.Time    `db:"start_at" validate:"required"`
	EndAt     sql.NullTime `db:"end_at"`
}

type Event struct {
	ID          int64           `db:"id"`
	Name        string          `db:"name" validate:"lte=128"`
	Description sql.NullString  `db:"description"`
	Image       *graphql.Upload `db:"image"`
	CreatedAt   time.Time       `db:"created_at"`
	StartAt     time.Time       `db:"start_at" validate:"required"`
	EndAt       sql.NullTime    `db:"end_at"`
}

type UpdateEvent struct {
	ID          int64 `validate:"required"`
	Name        sql.NullString
	Description sql.NullString
	Image       *graphql.Upload
	StartAt     sql.NullTime
	EndAt       sql.NullTime
}
