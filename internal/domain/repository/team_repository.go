package repository

import (
	"champi-maker/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type TeamRepository interface {
	Create(ctx context.Context, team *entity.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Team, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entity.Team, error)
}
