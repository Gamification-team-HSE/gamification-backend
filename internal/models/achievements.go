package models

import (
	"database/sql"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type RepoAchievement struct {
	ID          int            `db:"id"`
	Name        string         `db:"string"`
	Description sql.NullString `db:"description"`
	Image       sql.NullString `db:"image"`
	EndAt       sql.NullTime   `db:"end_at"`
	CreatedAt   time.Time      `db:"created_at"`
	Rules       *Rules         `db:"rules"`
}

type Rules struct {
	Blocks []*RulesBlock
}

type RulesBlock struct {
	EventRules         []*EventRule       `json:"event_rules"`
	StatRules          []*StatRule        `json:"stat_rules"`
	ConnectionOperator ConnectionOperator `json:"connection_operator"`
}

type EventRule struct {
	EventID         int  `json:"event_id"`
	NeedParticipate bool `json:"need_participate"`
}

type StatRule struct {
	StatID      int        `json:"stat_id"`
	TargetValue int        `json:"target_value"`
	Comparison  Comparison `json:"comparison"`
}

type Comparison string

const (
	ComparisonInvalid     Comparison = "invalid_comparison"
	ComparisonGreaterThan Comparison = ">"
	ComparisonLesserThan  Comparison = "<"
	ComparisonEquals      Comparison = "="
	ComparisonNotEquals   Comparison = "!="
)

type ConnectionOperator string

const (
	ConnectionOperatorInvalid = "invalid_operator"
	ConnectionOperatorAnd     = "and"
	ConnectionOperatorOr      = "or"
)

type CreateAchievement struct {
	Name        string
	Description string
	Image       *graphql.Upload
	Rules       *Rules
	EndAt       time.Time
}

type UpdateAchievement struct {
	ID          int
	Name        string
	Description string
	Image       *graphql.Upload
	Rules       *Rules
	EndAt       time.Time
}
