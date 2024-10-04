package repository

import (
	"champi-maker/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type ChampionshipRepository interface {
	Create(ctx context.Context, championship *entity.Championship) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Championship, error)
	Update(ctx context.Context, championship *entity.Championship) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entity.Championship, error)
}
