package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Image       sql.NullString `db:"image"`
	CreatedAt   time.Time      `db:"created_at"`
	StartAt     time.Time      `db:"start_at"`
	EndAt       sql.NullTime   `db:"end_at"`
}
