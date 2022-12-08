package models

import (
	"database/sql"
	"time"
)

type Stat struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	StartAt     time.Time      `db:"start_at"`
	Period      string         `db:"period"`
	SeqPeriod   sql.NullString `db:"seq_period"`
}
