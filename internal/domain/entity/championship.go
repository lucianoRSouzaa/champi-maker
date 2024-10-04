package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ChampionshipType string
type TiebreakerMethod string
type ProgressionType string

const (
	ChampionshipTypeLeague ChampionshipType = "league"
	ChampionshipTypeCup    ChampionshipType = "cup"

	TiebreakerPenalties TiebreakerMethod = "penalties"
	TiebreakerExtraTime TiebreakerMethod = "extra_time"

	ProgressionFixed      ProgressionType = "fixed"
	ProgressionRandomDraw ProgressionType = "random_draw"
)

type Championship struct {
	ID               uuid.UUID        `json:"id" validate:"required"`
	Name             string           `json:"name" validate:"required,min=2,max=100"`
	Type             ChampionshipType `json:"type" validate:"required,oneof=league cup"`
	TiebreakerMethod TiebreakerMethod `json:"tiebreaker_method" validate:"required,oneof=penalties extra_time"`
	ProgressionType  ProgressionType  `json:"progression_type" validate:"required,oneof=fixed random_draw"`
	CreatedAt        time.Time        `json:"created_at" validate:"required"`
	UpdatedAt        time.Time        `json:"updated_at" validate:"required"`
}

func (c *Championship) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
