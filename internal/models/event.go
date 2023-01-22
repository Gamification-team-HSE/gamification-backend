package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name" validate:"lte=128"`
	Description sql.NullString `db:"description" validate:"lte=1024"`
	Image       sql.NullString `db:"image"`
	CreatedAt   time.Time      `db:"created_at" `
	StartAt     time.Time      `db:"start_at" validate:"gtefield created_at"`
	EndAt       sql.NullTime   `db:"end_at" validate:"gtfield start_at"`
}
