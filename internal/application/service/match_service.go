package service

import (
	"champi-maker/internal/application"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"errors"
	"math"
	"math/rand/v2"
	"time"

	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MatchResultUpdate struct {
	ScoreHome          int
	ScoreAway          int
	HasExtraTime       bool
	ScoreHomeExtraTime int
	ScoreAwayExtraTime int
	HasPenalties       bool
	ScoreHomePenalties int
	ScoreAwayPenalties int
}

type MatchService interface {
	GenerateMatches(ctx context.Context, message application.ChampionshipCreatedMessage) error
	UpdateMatchResult(ctx context.Context, matchID uuid.UUID, result MatchResultUpdate) error
	GetMatchByID(ctx context.Context, matchID uuid.UUID) (*entity.Match, error)
	ListMatchesByChampionship(ctx context.Context, championshipID uuid.UUID) ([]*entity.Match, error)
}

type matchService struct {
	matchRepo         repository.MatchRepository
	championshipRepo  repository.ChampionshipRepository
	teamRepo          repository.TeamRepository
	statisticsService StatisticsService
}

func NewMatchService(
	matchRepo repository.MatchRepository,
	championshipRepo repository.ChampionshipRepository,
	teamRepo repository.TeamRepository,
	statisticsService StatisticsService,
) MatchService {
	return &matchService{
		matchRepo:         matchRepo,
		championshipRepo:  championshipRepo,
		teamRepo:          teamRepo,
		statisticsService: statisticsService,
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

func (s *matchService) UpdateMatchResult(ctx context.Context, matchID uuid.UUID, result MatchResultUpdate) error {
	// Iniciar transação
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

	match, err := s.matchRepo.GetByIDWithTx(ctx, tx, matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return fmt.Errorf("partida com ID %s não encontrada", matchID)
	}

	// Atualizar o resultado da partida
	match.ScoreHome = result.ScoreHome
	match.ScoreAway = result.ScoreAway
	match.HasExtraTime = result.HasExtraTime
	if result.HasExtraTime {
		match.ScoreHomeExtraTime = result.ScoreHomeExtraTime
		match.ScoreAwayExtraTime = result.ScoreAwayExtraTime
	}
	match.HasPenalties = result.HasPenalties
	if result.HasPenalties {
		match.ScoreHomePenalties = result.ScoreHomePenalties
		match.ScoreAwayPenalties = result.ScoreAwayPenalties
	}
	match.Status = entity.MatchStatusFinished
	match.UpdatedAt = time.Now()

	// Determinar o time vencedor
	winnerTeamID, err := s.determineWinner(match)
	if err != nil {
		return err
	}
	match.WinnerTeamID = winnerTeamID

	// Atualizar a partida no banco de dados
	if err := s.matchRepo.UpdateWithTx(ctx, tx, match); err != nil {
		return err
	}

	// Propagar o vencedor para a próxima fase (apenas para campeonatos do tipo Copa)
	championship, err := s.championshipRepo.GetByID(ctx, match.ChampionshipID)
	if err != nil {
		return err
	}
	if championship.Type == entity.ChampionshipTypeCup && match.ParentMatchID != nil {
		err = s.propagateWinner(ctx, tx, *match.ParentMatchID, *winnerTeamID)
		if err != nil {
			return err
		}
	}

	// Commit da transação de partida
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Atualizar estatísticas (fora da transação anterior)
	if championship.Type == entity.ChampionshipTypeLeague {
		if err := s.statisticsService.UpdateStatisticsAfterMatch(ctx, match); err != nil {
			return err
		}
	}

	return nil
}

func (s *matchService) determineWinner(match *entity.Match) (*uuid.UUID, error) {
	var homeGoals, awayGoals int

	homeGoals = match.ScoreHome
	awayGoals = match.ScoreAway

	if match.HasExtraTime {
		homeGoals += match.ScoreHomeExtraTime
		awayGoals += match.ScoreAwayExtraTime
	}

	if homeGoals > awayGoals {
		return match.HomeTeamID, nil
	} else if awayGoals > homeGoals {
		return match.AwayTeamID, nil
	} else {
		if match.HasPenalties {
			if match.ScoreHomePenalties > match.ScoreAwayPenalties {
				return match.HomeTeamID, nil
			} else if match.ScoreAwayPenalties > match.ScoreHomePenalties {
				return match.AwayTeamID, nil
			} else {
				return nil, errors.New("empate nas penalidades não é permitido")
			}

		} else {
			return nil, errors.New("partida terminou empatada e não há penalidades")
		}
	}
}

func (s *matchService) propagateWinner(ctx context.Context, tx pgx.Tx, parentMatchID uuid.UUID, winnerTeamID uuid.UUID) error {
	parentMatch, err := s.matchRepo.GetByIDWithTx(ctx, tx, parentMatchID)
	if err != nil {
		return err
	}
	if parentMatch == nil {
		return fmt.Errorf("partida pai com ID %s não encontrada", parentMatchID)
	}

	// Determinar se o vencedor vai para o lado esquerdo ou direito
	if parentMatch.LeftChildMatchID != nil && *parentMatch.LeftChildMatchID == parentMatchID {
		parentMatch.HomeTeamID = &winnerTeamID
	} else if parentMatch.RightChildMatchID != nil && *parentMatch.RightChildMatchID == parentMatchID {
		parentMatch.AwayTeamID = &winnerTeamID
	} else {
		return fmt.Errorf("partida com ID %s não é filha da partida com ID %s", parentMatchID, parentMatch.ID)
	}

	parentMatch.UpdatedAt = time.Now()

	// Atualizar a partida pai
	if err := s.matchRepo.UpdateWithTx(ctx, tx, parentMatch); err != nil {
		return err
	}

	return nil
}

func (s *matchService) GetMatchByID(ctx context.Context, matchID uuid.UUID) (*entity.Match, error) {
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, fmt.Errorf("partida com ID %s não encontrada", matchID)
	}
	return match, nil
}

func (s *matchService) ListMatchesByChampionship(ctx context.Context, championshipID uuid.UUID) ([]*entity.Match, error) {
	matches, err := s.matchRepo.GetByChampionshipID(ctx, championshipID)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
