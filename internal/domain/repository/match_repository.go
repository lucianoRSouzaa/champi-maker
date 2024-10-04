package repository

import (
	"champi-maker/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type MatchRepository interface {
	Create(ctx context.Context, match *entity.Match) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	Update(ctx context.Context, match *entity.Match) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByChampionshipID(ctx context.Context, championshipID uuid.UUID) ([]*entity.Match, error)
	GetByPhase(ctx context.Context, championshipID uuid.UUID, phase int) ([]*entity.Match, error)
}
