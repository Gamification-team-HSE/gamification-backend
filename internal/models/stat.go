package models

import (
	"database/sql"
	"time"
)

type Stat struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name" validate:"required"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	StartAt     time.Time      `db:"start_at"`
	Period      string         `db:"period"`
	SeqPeriod   sql.NullString `db:"seq_period"`
}

type GetStatsResponse struct {
	Stats []*Stat
	Total int
}

type UpdateStat struct {
	ID          int
	Name        string
	Description string
	StartedAt   time.Time
	Period      string
	SeqPeriod   string
}
