package models

import (
	"database/sql"
	"github.com/99designs/gqlgen/graphql"
	"time"
)

// DbEvent Select из базы возвращает событие в текущей структуре
type DbEvent struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name" validate:"lte=128"`
	Description sql.NullString `db:"description" validate:"lte=1024"`
	Image       string         `db:"image"`
	CreatedAt   time.Time      `db:"created_at"`
	StartAt     time.Time      `db:"start_at" validate:"required,gtefield=CreatedAt"`
	EndAt       sql.NullTime   `db:"end_at" validate:"gtfield=StartAt"`
}

type Event struct {
	ID          int64           `db:"id"`
	Name        string          `db:"name" validate:"lte=128"`
	Description sql.NullString  `db:"description" validate:"lte=1024"`
	Image       *graphql.Upload `db:"image"`
	CreatedAt   time.Time       `db:"created_at"`
	StartAt     time.Time       `db:"start_at" validate:"required"`
	EndAt       sql.NullTime    `db:"end_at"`
}

type UpdateEvent struct {
	ID          int64          `validate:"required"`
	Name        sql.NullString `validate:"lte=128"`
	Description sql.NullString `validate:"lte=1024"`
	Image       *graphql.Upload
	StartAt     sql.NullTime
	EndAt       sql.NullTime
}
