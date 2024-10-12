package port

import (
	"context"

	"github.com/google/uuid"
)

type MessagePublisher interface {
	PublishChampionshipCreated(ctx context.Context, championshipID uuid.UUID, teamIDs []uuid.UUID) error
}
