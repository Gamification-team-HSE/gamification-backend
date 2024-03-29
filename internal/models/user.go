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

type FullUser struct {
	*User
	Stats        []*UserStat
	Events       []*UserEvent
	Achievements []*UserAch
	Place        int
}

type UserStat struct {
	StatID      int            `db:"stat_id"`
	Value       int            `db:"value"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}

type UserEvent struct {
	EventID     int            `db:"event_id"`
	CreatedAt   time.Time      `db:"created_at"`
	Image       sql.NullString `db:"image"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}

type UserAch struct {
	AchID       int            `db:"ach_id"`
	CreatedAt   time.Time      `db:"created_at"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Image       sql.NullString `db:"image"`
}

type UserRatingByAch struct {
	UserID    int            `db:"user_id"`
	Name      sql.NullString `db:"name"`
	Email     string         `db:"email"`
	Place     int
	TotalAchs int            `db:"total_ach"`
	Avatar    sql.NullString `db:"avatar"`
}

type UserRatingByStat struct {
	UserID int            `db:"user_id"`
	Name   sql.NullString `db:"name"`
	Email  string         `db:"email"`
	Place  int
	Value  int            `db:"value"`
	Avatar sql.NullString `db:"avatar"`
}

type RatingByAchs struct {
	Total int
	Users []*UserRatingByAch
}

type RatingByStat struct {
	StatID int
	Total  int
	Users  []*UserRatingByStat
}
