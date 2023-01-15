package models

import (
	"database/sql"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type User struct {
	ID        int64          `db:"id"`
	ForeignID sql.NullString `db:"foreign_id"`
	Email     string         `db:"email" validate:"required,email"`
	CreatedAt time.Time      `db:"created_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	Role      Role           `db:"role" validate:"required"`
	Name      sql.NullString `db:"name" validate:"omitempty,max=256"`
	Avatar    sql.NullString `db:"avatar"`
}

type Role string

const (
	SuperAdminRole Role = "super_admin"
	AdminRole      Role = "admin"
	DefaultRole    Role = "user"
)

type UserFilter struct {
	Active bool
	Banned bool
	Admins bool
}

type UsersTotalInfo struct {
	Active int `db:"active"`
	Banned int `db:"banned"`
	Admins int `db:"admins"`
}

type GetUsersResponse struct {
	Users []*User
	Total *UsersTotalInfo
}

type UpdateUser struct {
	ID     int    `validate:"required"`
	Email  string `validate:"omitempty,email"`
	Name   string `validate:"omitempty,max=256"`
	Avatar *graphql.Upload
}
