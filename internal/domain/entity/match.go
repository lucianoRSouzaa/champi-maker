package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MatchStatus string

const (
	MatchStatusScheduled  MatchStatus = "scheduled"
	MatchStatusInProgress MatchStatus = "in_progress"
	MatchStatusFinished   MatchStatus = "finished"
)

type Match struct {
	ID                 uuid.UUID   `json:"id" validate:"required"`
	ChampionshipID     uuid.UUID   `json:"championship_id" validate:"required"`
	HomeTeamID         *uuid.UUID  `json:"home_team_id,omitempty" validate:"omitempty"`
	AwayTeamID         *uuid.UUID  `json:"away_team_id,omitempty" validate:"omitempty"`
	MatchDate          *time.Time  `json:"match_date,omitempty" validate:"omitempty"`
	Status             MatchStatus `json:"status" validate:"required,oneof=scheduled in_progress finished"`
	ScoreHome          int         `json:"score_home" validate:"gte=0"`
	ScoreAway          int         `json:"score_away" validate:"gte=0"`
	HasExtraTime       bool        `json:"has_extra_time"`
	ScoreHomeExtraTime int         `json:"score_home_extra_time" validate:"gte=0"`
	ScoreAwayExtraTime int         `json:"score_away_extra_time" validate:"gte=0"`
	HasPenalties       bool        `json:"has_penalties"`
	ScoreHomePenalties int         `json:"score_home_penalties" validate:"gte=0"`
	ScoreAwayPenalties int         `json:"score_away_penalties" validate:"gte=0"`
	WinnerTeamID       *uuid.UUID  `json:"winner_team_id,omitempty" validate:"omitempty"`
	Phase              int         `json:"phase" validate:"required,gte=1"`
	ParentMatchID      *uuid.UUID  `json:"parent_match_id,omitempty" validate:"omitempty"`
	LeftChildMatchID   *uuid.UUID  `json:"left_child_match_id,omitempty" validate:"omitempty"`
	RightChildMatchID  *uuid.UUID  `json:"right_child_match_id,omitempty" validate:"omitempty"`
	CreatedAt          time.Time   `json:"created_at" validate:"required"`
	UpdatedAt          time.Time   `json:"updated_at" validate:"required"`
}

func (m *Match) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}
