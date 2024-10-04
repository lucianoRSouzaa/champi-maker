package repository

import (
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type statisticsRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewStatisticsRepositoryPg(pool *pgxpool.Pool) repository.StatisticsRepository {
	return &statisticsRepositoryPg{pool: pool}
}

func (r *statisticsRepositoryPg) Create(ctx context.Context, stats *entity.Statistics) error {
	query := `
        INSERT INTO statistics (
            id, championship_id, team_id, matches_played, wins, draws, losses,
            goals_for, goals_against, goal_difference, points, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7,
                $8, $9, $10, $11, $12, $13)
    `
	_, err := r.pool.Exec(ctx, query,
		stats.ID,
		stats.ChampionshipID,
		stats.TeamID,
		stats.MatchesPlayed,
		stats.Wins,
		stats.Draws,
		stats.Losses,
		stats.GoalsFor,
		stats.GoalsAgainst,
		stats.GoalDifference,
		stats.Points,
		stats.CreatedAt,
		stats.UpdatedAt,
	)
	return err
}

func (r *statisticsRepositoryPg) GetByID(ctx context.Context, id uuid.UUID) (*entity.Statistics, error) {
	query := `
        SELECT
            id, championship_id, team_id, matches_played, wins, draws, losses,
            goals_for, goals_against, goal_difference, points, created_at, updated_at
        FROM statistics
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)

	var stats entity.Statistics
	err := row.Scan(
		&stats.ID,
		&stats.ChampionshipID,
		&stats.TeamID,
		&stats.MatchesPlayed,
		&stats.Wins,
		&stats.Draws,
		&stats.Losses,
		&stats.GoalsFor,
		&stats.GoalsAgainst,
		&stats.GoalDifference,
		&stats.Points,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Estatística não encontrada
		}
		return nil, err
	}

	return &stats, nil
}

func (r *statisticsRepositoryPg) GetByChampionshipAndTeam(ctx context.Context, championshipID, teamID uuid.UUID) (*entity.Statistics, error) {
	query := `
        SELECT
            id, championship_id, team_id, matches_played, wins, draws, losses,
            goals_for, goals_against, goal_difference, points, created_at, updated_at
        FROM statistics
        WHERE championship_id = $1 AND team_id = $2
    `
	row := r.pool.QueryRow(ctx, query, championshipID, teamID)

	var stats entity.Statistics
	err := row.Scan(
		&stats.ID,
		&stats.ChampionshipID,
		&stats.TeamID,
		&stats.MatchesPlayed,
		&stats.Wins,
		&stats.Draws,
		&stats.Losses,
		&stats.GoalsFor,
		&stats.GoalsAgainst,
		&stats.GoalDifference,
		&stats.Points,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Estatística não encontrada
		}
		return nil, err
	}

	return &stats, nil
}

func (r *statisticsRepositoryPg) Update(ctx context.Context, stats *entity.Statistics) error {
	query := `
        UPDATE statistics
        SET
            matches_played = $1,
            wins = $2,
            draws = $3,
            losses = $4,
            goals_for = $5,
            goals_against = $6,
            goal_difference = $7,
            points = $8,
            updated_at = $9
        WHERE id = $10
    `
	commandTag, err := r.pool.Exec(ctx, query,
		stats.MatchesPlayed,
		stats.Wins,
		stats.Draws,
		stats.Losses,
		stats.GoalsFor,
		stats.GoalsAgainst,
		stats.GoalDifference,
		stats.Points,
		time.Now(),
		stats.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *statisticsRepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM statistics
        WHERE id = $1
    `
	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were deleted")
	}

	return nil
}

func (r *statisticsRepositoryPg) ListByChampionship(ctx context.Context, championshipID uuid.UUID) ([]*entity.Statistics, error) {
	query := `
        SELECT
            id, championship_id, team_id, matches_played, wins, draws, losses,
            goals_for, goals_against, goal_difference, points, created_at, updated_at
        FROM statistics
        WHERE championship_id = $1
        ORDER BY points DESC, goal_difference DESC, goals_for DESC
    `
	rows, err := r.pool.Query(ctx, query, championshipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statsList []*entity.Statistics

	for rows.Next() {
		var stats entity.Statistics
		err := rows.Scan(
			&stats.ID,
			&stats.ChampionshipID,
			&stats.TeamID,
			&stats.MatchesPlayed,
			&stats.Wins,
			&stats.Draws,
			&stats.Losses,
			&stats.GoalsFor,
			&stats.GoalsAgainst,
			&stats.GoalDifference,
			&stats.Points,
			&stats.CreatedAt,
			&stats.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		statsList = append(statsList, &stats)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statsList, nil
}
