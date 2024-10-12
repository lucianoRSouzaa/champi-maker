package service

import (
	"champi-maker/internal/application"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"math"
	"math/rand/v2"
	"time"

	"context"
	"fmt"

	"github.com/google/uuid"
)

type MatchService interface {
	GenerateMatches(ctx context.Context, message application.ChampionshipCreatedMessage) error
	// Outros métodos conforme necessário
}

type matchService struct {
	matchRepo        repository.MatchRepository
	championshipRepo repository.ChampionshipRepository
	teamRepo         repository.TeamRepository
}

func NewMatchService(
	matchRepo repository.MatchRepository,
	championshipRepo repository.ChampionshipRepository,
	teamRepo repository.TeamRepository,
) MatchService {
	return &matchService{
		matchRepo:        matchRepo,
		championshipRepo: championshipRepo,
		teamRepo:         teamRepo,
	}
}

func (s *matchService) GenerateMatches(ctx context.Context, message application.ChampionshipCreatedMessage) error {
	// Iniciar a transação
	tx, err := s.matchRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Definir as restrições como deferidas
	_, err = tx.Exec(ctx, "SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		return err
	}

	championship, err := s.championshipRepo.GetByID(ctx, message.ChampionshipID)
	if err != nil {
		return err
	}
	if championship == nil {
		return fmt.Errorf("championship with ID %s not found", message.ChampionshipID)
	}

	// Gerar as partidas com base no tipo de campeonato
	var matches []*entity.Match

	switch championship.Type {
	case entity.ChampionshipTypeLeague:
		matches, err = s.generateLeagueMatches(ctx, championship, message.TeamIDs)
	case entity.ChampionshipTypeCup:
		matches, err = s.generateCupMatches(ctx, championship, message.TeamIDs)
	default:
		return fmt.Errorf("unknown championship type: %s", championship.Type)
	}

	if err != nil {
		return err
	}

	// Salvar as partidas no banco de dados
	for _, match := range matches {
		if err := s.matchRepo.CreateWithTx(ctx, tx, match); err != nil {
			return err
		}
	}

	return nil
}

func (s *matchService) generateLeagueMatches(ctx context.Context, championship *entity.Championship, teamIDs []uuid.UUID) ([]*entity.Match, error) {
	var matches []*entity.Match

	for i := 0; i < len(teamIDs); i++ {
		for j := i + 1; j < len(teamIDs); j++ {
			match := &entity.Match{
				ID:             uuid.New(),
				ChampionshipID: championship.ID,
				HomeTeamID:     &teamIDs[i],
				AwayTeamID:     &teamIDs[j],
				MatchDate:      nil, // Pode ser definido posteriormente
				Status:         entity.MatchStatusScheduled,
				Phase:          1,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			matches = append(matches, match)
		}
	}

	return matches, nil
}

func (s *matchService) generateCupMatches(ctx context.Context, championship *entity.Championship, teamIDs []uuid.UUID) ([]*entity.Match, error) {
	numTeams := len(teamIDs)
	numPhases := int(math.Ceil(math.Log2(float64(numTeams))))
	numSlots := int(math.Pow(2, float64(numPhases)))

	// Criar slots, preenchendo com nil para byes
	slots := make([]*uuid.UUID, numSlots)
	for i := 0; i < numSlots; i++ {
		if i < numTeams {
			slots[i] = &teamIDs[i]
		} else {
			slots[i] = nil // Bye
		}
	}

	// Embaralhar os slots se necessário
	if championship.ProgressionType == entity.ProgressionRandomDraw {
		rand.Shuffle(len(slots), func(i, j int) { slots[i], slots[j] = slots[j], slots[i] })
	}

	// Inicializar as fases
	phases := make([][]*entity.Match, numPhases)

	// Gerar as partidas da Fase 1 (fase inicial)
	phases[0] = make([]*entity.Match, 0)
	for i := 0; i < len(slots); i += 2 {
		homeTeamID := slots[i]
		awayTeamID := slots[i+1]
		match := &entity.Match{
			ID:             uuid.New(),
			ChampionshipID: championship.ID,
			HomeTeamID:     homeTeamID,
			AwayTeamID:     awayTeamID,
			MatchDate:      nil, // Pode ser definido posteriormente
			Status:         entity.MatchStatusScheduled,
			Phase:          1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		phases[0] = append(phases[0], match)
	}

	if championship.ProgressionType == entity.ProgressionFixed {
		// Gerar as partidas das fases subsequentes
		for phase := 1; phase < numPhases; phase++ {
			previousPhaseMatches := phases[phase-1]
			numMatches := len(previousPhaseMatches)
			currentPhaseMatches := make([]*entity.Match, 0)
			for i := 0; i < numMatches; i += 2 {
				leftChildMatch := previousPhaseMatches[i]
				var rightChildMatch *entity.Match
				if i+1 < numMatches {
					rightChildMatch = previousPhaseMatches[i+1]
				} else {
					rightChildMatch = nil // Bye
				}
				match := &entity.Match{
					ID:                uuid.New(),
					ChampionshipID:    championship.ID,
					LeftChildMatchID:  &leftChildMatch.ID,
					RightChildMatchID: nil,
					MatchDate:         nil,
					Status:            entity.MatchStatusScheduled,
					Phase:             phase + 1,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}
				if rightChildMatch != nil {
					match.RightChildMatchID = &rightChildMatch.ID
				}
				// Atualizar o ParentMatchID nas partidas filhas
				leftChildMatch.ParentMatchID = &match.ID
				if rightChildMatch != nil {
					rightChildMatch.ParentMatchID = &match.ID
				}
				currentPhaseMatches = append(currentPhaseMatches, match)
			}
			phases[phase] = currentPhaseMatches
		}
	}

	// Coletar todas as partidas em uma única lista
	matches := make([]*entity.Match, 0)
	for _, phaseMatches := range phases {
		matches = append(matches, phaseMatches...)
	}

	return matches, nil
}
