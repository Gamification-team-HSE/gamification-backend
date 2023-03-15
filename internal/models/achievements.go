package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type RepoAchievement struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Image       sql.NullString `db:"image"`
	EndAt       sql.NullTime   `db:"end_at"`
	CreatedAt   time.Time      `db:"created_at"`
	Rules       *Rules         `db:"rules"`
}

type Rules struct {
	Blocks []*RulesBlock
}

func (r *Rules) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to cast value ti ")
	}
	return json.Unmarshal(b, r)
}

func (r *Rules) Value() (driver.Value, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return b, nil
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

type GetAchievementsResponse struct {
	Achievements []*RepoAchievement
	Total        int
}
