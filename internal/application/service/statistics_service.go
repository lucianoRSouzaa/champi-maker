package service

import (
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StatisticsService interface {
	GenerateInitialStatistics(ctx context.Context, championshipID uuid.UUID, teamIDs []uuid.UUID) error
	UpdateStatisticsAfterMatch(ctx context.Context, match *entity.Match) error
	GetStatisticsByChampionship(ctx context.Context, championshipID uuid.UUID) ([]*entity.Statistics, error)
}

type statisticsService struct {
	statisticsRepo   repository.StatisticsRepository
	championshipRepo repository.ChampionshipRepository
	teamRepo         repository.TeamRepository
}

func NewStatisticsService(
	statisticsRepo repository.StatisticsRepository,
	championshipRepo repository.ChampionshipRepository,
	teamRepo repository.TeamRepository,
) StatisticsService {
	return &statisticsService{
		statisticsRepo:   statisticsRepo,
		championshipRepo: championshipRepo,
		teamRepo:         teamRepo,
	}
}

func (s *statisticsService) GenerateInitialStatistics(ctx context.Context, championshipID uuid.UUID, teamIDs []uuid.UUID) error {
	championship, err := s.championshipRepo.GetByID(ctx, championshipID)
	if err != nil {
		return err
	}
	if championship == nil {
		return fmt.Errorf("campeonato com ID %s não encontrado", championshipID)
	}
	if championship.Type != entity.ChampionshipTypeLeague {
		return fmt.Errorf("estatísticas iniciais só são geradas para campeonatos do tipo liga")
	}

	for _, teamID := range teamIDs {
		team, err := s.teamRepo.GetByID(ctx, teamID)
		if err != nil {
			return err
		}
		if team == nil {
			return fmt.Errorf("time com ID %s não encontrado", teamID)
		}

		stats := &entity.Statistics{
			ID:             uuid.New(),
			ChampionshipID: championshipID,
			TeamID:         teamID,
			MatchesPlayed:  0,
			Wins:           0,
			Draws:          0,
			Losses:         0,
			GoalsFor:       0,
			GoalsAgainst:   0,
			GoalDifference: 0,
			Points:         0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.statisticsRepo.Create(ctx, stats); err != nil {
			return err
		}
	}

	return nil
}

func (s *statisticsService) UpdateStatisticsAfterMatch(ctx context.Context, match *entity.Match) error {
	// Verificar se a partida está concluída
	if match.Status != entity.MatchStatusFinished {
		return fmt.Errorf("a partida com ID %s não está concluída", match.ID)
	}

	// Obter o campeonato
	championship, err := s.championshipRepo.GetByID(ctx, match.ChampionshipID)
	if err != nil {
		return err
	}
	if championship == nil {
		return fmt.Errorf("campeonato com ID %s não encontrado", match.ChampionshipID)
	}

	// Atualizar estatísticas apenas para campeonatos do tipo liga
	if championship.Type != entity.ChampionshipTypeLeague {
		return nil // Não faz nada se não for liga
	}

	// Iniciar transação
	tx, err := s.statisticsRepo.BeginTx(ctx)
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

	// Atualizar estatísticas do time da casa
	if err := s.updateTeamStatistics(ctx, tx, championship.ID, *match.HomeTeamID, match, true); err != nil {
		return err
	}

	// Atualizar estatísticas do time visitante
	if err := s.updateTeamStatistics(ctx, tx, championship.ID, *match.AwayTeamID, match, false); err != nil {
		return err
	}

	return nil
}

func (s *statisticsService) updateTeamStatistics(ctx context.Context, tx pgx.Tx, championshipID, teamID uuid.UUID, match *entity.Match, isHomeTeam bool) error {
	stats, err := s.statisticsRepo.GetByChampionshipAndTeamWithTx(ctx, tx, championshipID, teamID)
	if err != nil {
		return err
	}
	if stats == nil {
		return fmt.Errorf("estatísticas não encontradas para o time %s no campeonato %s", teamID, championshipID)
	}

	// Atualizar estatísticas básicas
	stats.MatchesPlayed += 1

	var goalsFor, goalsAgainst int
	if isHomeTeam {
		goalsFor = match.ScoreHome
		goalsAgainst = match.ScoreAway
	} else {
		goalsFor = match.ScoreAway
		goalsAgainst = match.ScoreHome
	}

	stats.GoalsFor += goalsFor
	stats.GoalsAgainst += goalsAgainst
	stats.GoalDifference = stats.GoalsFor - stats.GoalsAgainst

	// Determinar resultado
	if goalsFor > goalsAgainst {
		stats.Wins += 1
		stats.Points += 3
	} else if goalsFor == goalsAgainst {
		stats.Draws += 1
		stats.Points += 1
	} else {
		stats.Losses += 1
	}

	stats.UpdatedAt = time.Now()

	// Atualizar no banco de dados
	if err := s.statisticsRepo.UpdateWithTx(ctx, tx, stats); err != nil {
		return err
	}

	return nil
}

func (s *statisticsService) GetStatisticsByChampionship(ctx context.Context, championshipID uuid.UUID) ([]*entity.Statistics, error) {
	statsList, err := s.statisticsRepo.ListByChampionship(ctx, championshipID)
	if err != nil {
		return nil, err
	}
	return statsList, nil
}
