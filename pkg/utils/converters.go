package utils

import (
	"database/sql"
)

func SqlNullStringToString(sqlString sql.NullString) *string {
	if sqlString.Valid {
		return &sqlString.String
	}
	return nil
}

func SqlNullTimeToTime(sqlTime sql.NullTime) int {
	if sqlTime.Valid {
		return int(sqlTime.Time.Unix())
	}
	return 0
}
