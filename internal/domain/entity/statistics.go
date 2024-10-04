package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Statistics struct {
	ID             uuid.UUID `json:"id" validate:"required"`
	ChampionshipID uuid.UUID `json:"championship_id" validate:"required"`
	TeamID         uuid.UUID `json:"team_id" validate:"required"`
	MatchesPlayed  int       `json:"matches_played" validate:"gte=0"`
	Wins           int       `json:"wins" validate:"gte=0"`
	Draws          int       `json:"draws" validate:"gte=0"`
	Losses         int       `json:"losses" validate:"gte=0"`
	GoalsFor       int       `json:"goals_for" validate:"gte=0"`
	GoalsAgainst   int       `json:"goals_against" validate:"gte=0"`
	GoalDifference int       `json:"goal_difference" validate:"gte=0"`
	Points         int       `json:"points" validate:"gte=0"`
	CreatedAt      time.Time `json:"created_at" validate:"required"`
	UpdatedAt      time.Time `json:"updated_at" validate:"required"`
}

func (s *Statistics) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
