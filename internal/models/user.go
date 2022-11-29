package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `db:"id"`
	ForeignID sql.NullString `db:"foreign_id"`
	Email     string         `db:"email" validate:"required,email"`
	CreatedAt time.Time      `db:"created_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	Role      Role           `db:"role" validate:"required"`
	Name      sql.NullString `db:"name"`
	Avatar    sql.NullString `db:"avatar"`
}

type Role string

const (
	SuperAdminRole Role = "super_admin"
	AdminRole      Role = "admin"
	DefaultRole    Role = "user"
)
