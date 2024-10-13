package repository

import (
	"champi-maker/internal/domain/entity"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MatchRepository interface {
	Create(ctx context.Context, match *entity.Match) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	Update(ctx context.Context, match *entity.Match) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByChampionshipID(ctx context.Context, championshipID uuid.UUID) ([]*entity.Match, error)
	GetByPhase(ctx context.Context, championshipID uuid.UUID, phase int) ([]*entity.Match, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, match *entity.Match) error
	GetByIDWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*entity.Match, error)
	UpdateWithTx(ctx context.Context, tx pgx.Tx, match *entity.Match) error
}
