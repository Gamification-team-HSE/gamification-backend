// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type Achievement struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	Rules       *Rules  `json:"rules"`
	EndAt       *int    `json:"end_at"`
	CreatedAt   int     `json:"created_at"`
}

type CreateAchievement struct {
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	Rules       *InputRules     `json:"rules"`
	EndAt       *int            `json:"end_at"`
}

type Event struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	CreatedAt   int             `json:"created_at"`
	StartAt     int             `json:"start_at"`
	EndAt       *int            `json:"end_at"`
}

type EventRule struct {
	EventID         int  `json:"event_id"`
	NeedParticipate bool `json:"need_participate"`
}

type FullUser struct {
	User         *User        `json:"user"`
	Stats        []*UserStat  `json:"stats"`
	Events       []*UserEvent `json:"events"`
	Achievements []*UserAch   `json:"achievements"`
}

type GetAchievementsResponse struct {
	Total        int            `json:"total"`
	Achievements []*Achievement `json:"achievements"`
}

type GetEvent struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	CreatedAt   int     `json:"created_at"`
	StartAt     int     `json:"start_at"`
	EndAt       *int    `json:"end_at"`
}

type GetEventsResponse struct {
	Total  int         `json:"total"`
	Events []*GetEvent `json:"events"`
}

type GetStatsResponse struct {
	Total int     `json:"total"`
	Stats []*Stat `json:"stats"`
}

type GetUsersResponse struct {
	Users []*User         `json:"users"`
	Total *UsersTotalInfo `json:"total"`
}

type InputEventRule struct {
	EventID         int  `json:"event_id"`
	NeedParticipate bool `json:"need_participate"`
}

type InputRuleBlock struct {
	EventsRules        []*InputEventRule   `json:"eventsRules"`
	StatRules          []*InputStatRule    `json:"statRules"`
	ConnectionOperator *ConnectionOperator `json:"connection_operator"`
}

type InputRules struct {
	Blocks []*InputRuleBlock `json:"blocks"`
}

type InputStatRule struct {
	StatID         int        `json:"stat_id"`
	TargetValue    int        `json:"target_value"`
	ComparisonType Comparison `json:"comparison_type"`
}

type NewEvent struct {
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	StartAt     int             `json:"start_at"`
	EndAt       *int            `json:"end_at"`
}

type NewStat struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StartAt     int     `json:"start_at"`
	Period      string  `json:"period"`
	SeqPeriod   *string `json:"seq_period"`
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

type RuleBlock struct {
	EventsRules        []*EventRule        `json:"eventsRules"`
	StatRules          []*StatRule         `json:"statRules"`
	ConnectionOperator *ConnectionOperator `json:"connection_operator"`
}

type Rules struct {
	Blocks []*RuleBlock `json:"blocks"`
}

type Stat struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   int     `json:"created_at"`
	StartAt     int     `json:"start_at"`
	Period      string  `json:"period"`
	SeqPeriod   *string `json:"seq_period"`
}

type StatRule struct {
	StatID         int        `json:"stat_id"`
	TargetValue    int        `json:"target_value"`
	ComparisonType Comparison `json:"comparison_type"`
}

type UpdateAchievement struct {
	ID          int             `json:"id"`
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	Rules       *InputRules     `json:"rules"`
	EndAt       *int            `json:"end_at"`
}

type UpdateEvent struct {
	ID          int             `json:"id"`
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Image       *graphql.Upload `json:"image"`
	StartAt     *int            `json:"start_at"`
	EndAt       *int            `json:"end_at"`
}

type UpdateStat struct {
	ID          int     `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	StartAt     *int    `json:"start_at"`
	Period      *string `json:"period"`
	SeqPeriod   *string `json:"seq_period"`
}

type UpdateUser struct {
	ID     int             `json:"id"`
	Email  *string         `json:"email"`
	Avatar *graphql.Upload `json:"avatar"`
	Name   *string         `json:"name"`
}

type User struct {
	ID        int     `json:"id"`
	ForeignID *string `json:"foreign_id"`
	Email     string  `json:"email"`
	CreatedAt int     `json:"created_at"`
	DeletedAt *int    `json:"deleted_at"`
	Role      Role    `json:"role"`
	Avatar    *string `json:"avatar"`
	Name      *string `json:"name"`
}

type UserAch struct {
	AchID       int     `json:"ach_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   int     `json:"created_at"`
	Image       *string `json:"image"`
}

type UserEvent struct {
	EventID     int     `json:"event_id"`
	Name        string  `json:"name"`
	Image       *string `json:"image"`
	Description *string `json:"description"`
	CreatedAt   int     `json:"created_at"`
}

type UserFilter struct {
	Active *bool `json:"active"`
	Banned *bool `json:"banned"`
	Admins *bool `json:"admins"`
}

type UserStat struct {
	StatID      int     `json:"stat_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Value       int     `json:"value"`
}

type UsersTotalInfo struct {
	Admins int `json:"admins"`
	Banned int `json:"banned"`
	Active int `json:"active"`
}

type Comparison string

const (
	ComparisonInvalidComparison Comparison = "InvalidComparison"
	ComparisonGreaterThan       Comparison = "GreaterThan"
	ComparisonEquals            Comparison = "Equals"
	ComparisonNotEquals         Comparison = "NotEquals"
	ComparisonLesserThan        Comparison = "LesserThan"
)

var AllComparison = []Comparison{
	ComparisonInvalidComparison,
	ComparisonGreaterThan,
	ComparisonEquals,
	ComparisonNotEquals,
	ComparisonLesserThan,
}

func (e Comparison) IsValid() bool {
	switch e {
	case ComparisonInvalidComparison, ComparisonGreaterThan, ComparisonEquals, ComparisonNotEquals, ComparisonLesserThan:
		return true
	}
	return false
}

func (e Comparison) String() string {
	return string(e)
}

func (e *Comparison) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Comparison(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Comparison", str)
	}
	return nil
}

func (e Comparison) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ConnectionOperator string

const (
	ConnectionOperatorInvalidConnectionOperator ConnectionOperator = "InvalidConnectionOperator"
	ConnectionOperatorAnd                       ConnectionOperator = "And"
	ConnectionOperatorOr                        ConnectionOperator = "Or"
)

var AllConnectionOperator = []ConnectionOperator{
	ConnectionOperatorInvalidConnectionOperator,
	ConnectionOperatorAnd,
	ConnectionOperatorOr,
}

func (e ConnectionOperator) IsValid() bool {
	switch e {
	case ConnectionOperatorInvalidConnectionOperator, ConnectionOperatorAnd, ConnectionOperatorOr:
		return true
	}
	return false
}

func (e ConnectionOperator) String() string {
	return string(e)
}

func (e *ConnectionOperator) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ConnectionOperator(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ConnectionOperator", str)
	}
	return nil
}

func (e ConnectionOperator) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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
