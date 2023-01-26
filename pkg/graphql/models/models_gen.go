// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Event struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	CreatedAt   time.Time       `json:"created_at"`
	StartAt     time.Time       `json:"start_at"`
	EndAt       *time.Time      `json:"end_at"`
}

type GetUsersResponse struct {
	Users []*User         `json:"users"`
	Total *UsersTotalInfo `json:"total"`
}

type NewEvent struct {
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	StartAt     time.Time       `json:"start_at"`
	EndAt       *time.Time      `json:"end_at"`
}

type NewStat struct {
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	StartAt     time.Time `json:"start_at"`
	Period      string    `json:"period"`
	SeqPeriod   *string   `json:"seq_period"`
}

type NewUser struct {
	ForeignID *string `json:"foreign_id"`
	Email     string  `json:"email"`
	Role      Role    `json:"role"`
	Name      *string `json:"name"`
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type Stat struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	StartAt     time.Time `json:"start_at"`
	Period      string    `json:"period"`
	SeqPeriod   *string   `json:"seq_period"`
}

type UpdateEvent struct {
	ID          int             `json:"id"`
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	StartAt     *time.Time      `json:"start_at"`
	EndAt       *time.Time      `json:"end_at"`
}

type UpdateUser struct {
	ID     int             `json:"id"`
	Email  *string         `json:"email"`
	Avatar *graphql.Upload `json:"avatar"`
	Name   *string         `json:"name"`
}

type User struct {
	ID        int        `json:"id"`
	ForeignID *string    `json:"foreign_id"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Role      Role       `json:"role"`
	Avatar    *string    `json:"avatar"`
	Name      *string    `json:"name"`
}

type UserFilter struct {
	Active *bool `json:"active"`
	Banned *bool `json:"banned"`
	Admins *bool `json:"admins"`
}

type UsersTotalInfo struct {
	Admins int `json:"admins"`
	Banned int `json:"banned"`
	Active int `json:"active"`
}

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleUser       Role = "user"
	RoleSuperAdmin Role = "super_admin"
)

var AllRole = []Role{
	RoleAdmin,
	RoleUser,
	RoleSuperAdmin,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleUser, RoleSuperAdmin:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
