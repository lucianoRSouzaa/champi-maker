package service

import (
	"champi-maker/internal/application/port"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChampionshipService interface {
	CreateChampionship(ctx context.Context, championship *entity.Championship, teamIDs []uuid.UUID) error
	// Outros métodos conforme necessário
}

type championshipService struct {
	championshipRepo repository.ChampionshipRepository
	teamRepo         repository.TeamRepository
	messagePublisher port.MessagePublisher
}

func NewChampionshipService(
	championshipRepo repository.ChampionshipRepository,
	teamRepo repository.TeamRepository,
	messagePublisher port.MessagePublisher,
) ChampionshipService {
	return &championshipService{
		championshipRepo: championshipRepo,
		teamRepo:         teamRepo,
		messagePublisher: messagePublisher,
	}
}

func (s *championshipService) CreateChampionship(ctx context.Context, championship *entity.Championship, teamIDs []uuid.UUID) error {
	if err := championship.Validate(); err != nil {
		return err
	}

	// Verificar se os times existem
	for _, teamID := range teamIDs {
		team, err := s.teamRepo.GetByID(ctx, teamID)
		if err != nil {
			return err
		}
		if team == nil {
			return fmt.Errorf("team with ID %s not found", teamID)
		}
	}

	// Definir IDs e timestamps
	championship.ID = uuid.New()
	championship.CreatedAt = time.Now()
	championship.UpdatedAt = time.Now()

	// Salvar o campeonato no banco de dados
	if err := s.championshipRepo.Create(ctx, championship); err != nil {
		return err
	}

	// Publicar mensagem usando a interface
	if err := s.messagePublisher.PublishChampionshipCreated(ctx, championship.ID, teamIDs); err != nil {
		return err
	}

	return nil
}
