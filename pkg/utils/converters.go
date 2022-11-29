package utils

import (
	"database/sql"
	"time"
)

func SqlNullStringToString(sqlString sql.NullString) *string {
	if sqlString.Valid {
		return &sqlString.String
	}
	return nil
}

func SqlNullTimeToTime(sqlTime sql.NullTime) *time.Time {
	if sqlTime.Valid {
		return &sqlTime.Time
	}
	return nil
}
