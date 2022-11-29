package utils

import (
	sq "github.com/Masterminds/squirrel"
)

func PgQB() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
