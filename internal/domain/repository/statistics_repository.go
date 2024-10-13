package repository

import (
	"champi-maker/internal/domain/entity"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StatisticsRepository interface {
	Create(ctx context.Context, stats *entity.Statistics) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Statistics, error)
	GetByChampionshipAndTeam(ctx context.Context, championshipID, teamID uuid.UUID) (*entity.Statistics, error)
	Update(ctx context.Context, stats *entity.Statistics) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByChampionship(ctx context.Context, championshipID uuid.UUID) ([]*entity.Statistics, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, stats *entity.Statistics) error
	GetByChampionshipAndTeamWithTx(ctx context.Context, tx pgx.Tx, championshipID, teamID uuid.UUID) (*entity.Statistics, error)
	UpdateWithTx(ctx context.Context, tx pgx.Tx, stats *entity.Statistics) error
}
