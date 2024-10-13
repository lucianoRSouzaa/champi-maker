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

type matchRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewMatchRepositoryPg(pool *pgxpool.Pool) repository.MatchRepository {
	return &matchRepositoryPg{pool: pool}
}

func (r *matchRepositoryPg) Create(ctx context.Context, match *entity.Match) error {
	query := `
        INSERT INTO matches (
            id, championship_id, home_team_id, away_team_id, match_date, status,
            score_home, score_away, has_extra_time, score_home_extra_time,
            score_away_extra_time, has_penalties, score_home_penalties,
            score_away_penalties, winner_team_id, phase, parent_match_id,
            left_child_match_id, right_child_match_id, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10,
            $11, $12, $13,
            $14, $15, $16, $17,
            $18, $19, $20, $21
        )
    `
	_, err := r.pool.Exec(ctx, query,
		match.ID,
		match.ChampionshipID,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.Status,
		match.ScoreHome,
		match.ScoreAway,
		match.HasExtraTime,
		match.ScoreHomeExtraTime,
		match.ScoreAwayExtraTime,
		match.HasPenalties,
		match.ScoreHomePenalties,
		match.ScoreAwayPenalties,
		match.WinnerTeamID,
		match.Phase,
		match.ParentMatchID,
		match.LeftChildMatchID,
		match.RightChildMatchID,
		match.CreatedAt,
		match.UpdatedAt,
	)
	return err
}

func (r *matchRepositoryPg) GetByID(ctx context.Context, id uuid.UUID) (*entity.Match, error) {
	query := `
        SELECT
            id, championship_id, home_team_id, away_team_id, match_date, status,
            score_home, score_away, has_extra_time, score_home_extra_time,
            score_away_extra_time, has_penalties, score_home_penalties,
            score_away_penalties, winner_team_id, phase, parent_match_id,
            left_child_match_id, right_child_match_id, created_at, updated_at
        FROM matches
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)

	var match entity.Match
	err := row.Scan(
		&match.ID,
		&match.ChampionshipID,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.MatchDate,
		&match.Status,
		&match.ScoreHome,
		&match.ScoreAway,
		&match.HasExtraTime,
		&match.ScoreHomeExtraTime,
		&match.ScoreAwayExtraTime,
		&match.HasPenalties,
		&match.ScoreHomePenalties,
		&match.ScoreAwayPenalties,
		&match.WinnerTeamID,
		&match.Phase,
		&match.ParentMatchID,
		&match.LeftChildMatchID,
		&match.RightChildMatchID,
		&match.CreatedAt,
		&match.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Match não encontrado
		}
		return nil, err
	}

	return &match, nil
}

func (r *matchRepositoryPg) Update(ctx context.Context, match *entity.Match) error {
	query := `
        UPDATE matches SET
            home_team_id = $1,
            away_team_id = $2,
            match_date = $3,
            status = $4,
            score_home = $5,
            score_away = $6,
            has_extra_time = $7,
            score_home_extra_time = $8,
            score_away_extra_time = $9,
            has_penalties = $10,
            score_home_penalties = $11,
            score_away_penalties = $12,
            winner_team_id = $13,
            phase = $14,
            parent_match_id = $15,
            left_child_match_id = $16,
            right_child_match_id = $17,
            updated_at = $18
        WHERE id = $19
    `
	commandTag, err := r.pool.Exec(ctx, query,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.Status,
		match.ScoreHome,
		match.ScoreAway,
		match.HasExtraTime,
		match.ScoreHomeExtraTime,
		match.ScoreAwayExtraTime,
		match.HasPenalties,
		match.ScoreHomePenalties,
		match.ScoreAwayPenalties,
		match.WinnerTeamID,
		match.Phase,
		match.ParentMatchID,
		match.LeftChildMatchID,
		match.RightChildMatchID,
		time.Now(),
		match.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *matchRepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM matches
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

func (r *matchRepositoryPg) GetByChampionshipID(ctx context.Context, championshipID uuid.UUID) ([]*entity.Match, error) {
	query := `
        SELECT
            id, championship_id, home_team_id, away_team_id, match_date, status,
            score_home, score_away, has_extra_time, score_home_extra_time,
            score_away_extra_time, has_penalties, score_home_penalties,
            score_away_penalties, winner_team_id, phase, parent_match_id,
            left_child_match_id, right_child_match_id, created_at, updated_at
        FROM matches
        WHERE championship_id = $1
        ORDER BY phase ASC, match_date ASC
    `
	rows, err := r.pool.Query(ctx, query, championshipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*entity.Match
	for rows.Next() {
		var match entity.Match
		err := rows.Scan(
			&match.ID,
			&match.ChampionshipID,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.MatchDate,
			&match.Status,
			&match.ScoreHome,
			&match.ScoreAway,
			&match.HasExtraTime,
			&match.ScoreHomeExtraTime,
			&match.ScoreAwayExtraTime,
			&match.HasPenalties,
			&match.ScoreHomePenalties,
			&match.ScoreAwayPenalties,
			&match.WinnerTeamID,
			&match.Phase,
			&match.ParentMatchID,
			&match.LeftChildMatchID,
			&match.RightChildMatchID,
			&match.CreatedAt,
			&match.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, &match)
	}

	return matches, nil
}

func (r *matchRepositoryPg) GetByPhase(ctx context.Context, championshipID uuid.UUID, phase int) ([]*entity.Match, error) {
	query := `
        SELECT
            id, championship_id, home_team_id, away_team_id, match_date, status,
            score_home, score_away, has_extra_time, score_home_extra_time,
            score_away_extra_time, has_penalties, score_home_penalties,
            score_away_penalties, winner_team_id, phase, parent_match_id,
            left_child_match_id, right_child_match_id, created_at, updated_at
        FROM matches
        WHERE championship_id = $1 AND phase = $2
        ORDER BY match_date ASC
    `
	rows, err := r.pool.Query(ctx, query, championshipID, phase)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*entity.Match
	for rows.Next() {
		var match entity.Match
		err := rows.Scan(
			&match.ID,
			&match.ChampionshipID,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.MatchDate,
			&match.Status,
			&match.ScoreHome,
			&match.ScoreAway,
			&match.HasExtraTime,
			&match.ScoreHomeExtraTime,
			&match.ScoreAwayExtraTime,
			&match.HasPenalties,
			&match.ScoreHomePenalties,
			&match.ScoreAwayPenalties,
			&match.WinnerTeamID,
			&match.Phase,
			&match.ParentMatchID,
			&match.LeftChildMatchID,
			&match.RightChildMatchID,
			&match.CreatedAt,
			&match.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, &match)
	}

	return matches, nil
}

func (r *matchRepositoryPg) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *matchRepositoryPg) CreateWithTx(ctx context.Context, tx pgx.Tx, match *entity.Match) error {
	query := `
        INSERT INTO matches (
            id, championship_id, home_team_id, away_team_id, match_date, status,
            score_home, score_away, has_extra_time, score_home_extra_time, score_away_extra_time,
            has_penalties, score_home_penalties, score_away_penalties, winner_team_id,
            created_at, updated_at, phase, parent_match_id, left_child_match_id, right_child_match_id
        )
        VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, $11,
            $12, $13, $14, $15,
            $16, $17, $18, $19, $20, $21
        )
    `
	_, err := tx.Exec(ctx, query,
		match.ID,
		match.ChampionshipID,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.Status,
		match.ScoreHome,
		match.ScoreAway,
		match.HasExtraTime,
		match.ScoreHomeExtraTime,
		match.ScoreAwayExtraTime,
		match.HasPenalties,
		match.ScoreHomePenalties,
		match.ScoreAwayPenalties,
		match.WinnerTeamID,
		match.CreatedAt,
		match.UpdatedAt,
		match.Phase,
		match.ParentMatchID,
		match.LeftChildMatchID,
		match.RightChildMatchID,
	)
	return err
}

func (r *matchRepositoryPg) GetByIDWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*entity.Match, error) {
	query := `
		SELECT
			id, championship_id, home_team_id, away_team_id, match_date, status,
			score_home, score_away, has_extra_time, score_home_extra_time,
			score_away_extra_time, has_penalties, score_home_penalties,
			score_away_penalties, winner_team_id, phase, parent_match_id,
			left_child_match_id, right_child_match_id, created_at, updated_at
		FROM matches
		WHERE id = $1
	`
	row := tx.QueryRow(ctx, query, id)

	var match entity.Match
	err := row.Scan(
		&match.ID,
		&match.ChampionshipID,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.MatchDate,
		&match.Status,
		&match.ScoreHome,
		&match.ScoreAway,
		&match.HasExtraTime,
		&match.ScoreHomeExtraTime,
		&match.ScoreAwayExtraTime,
		&match.HasPenalties,
		&match.ScoreHomePenalties,
		&match.ScoreAwayPenalties,
		&match.WinnerTeamID,
		&match.Phase,
		&match.ParentMatchID,
		&match.LeftChildMatchID,
		&match.RightChildMatchID,
		&match.CreatedAt,
		&match.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

func (r *matchRepositoryPg) UpdateWithTx(ctx context.Context, tx pgx.Tx, match *entity.Match) error {
	query := `
		UPDATE matches SET
			home_team_id = $1,
			away_team_id = $2,
			match_date = $3,
			status = $4,
			score_home = $5,
			score_away = $6,
			has_extra_time = $7,
			score_home_extra_time = $8,
			score_away_extra_time = $9,
			has_penalties = $10,
			score_home_penalties = $11,
			score_away_penalties = $12,
			winner_team_id = $13,
			phase = $14,
			parent_match_id = $15,
			left_child_match_id = $16,
			right_child_match_id = $17,
			updated_at = $18
		WHERE id = $19
	`
	commandTag, err := tx.Exec(ctx, query,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.Status,
		match.ScoreHome,
		match.ScoreAway,
		match.HasExtraTime,
		match.ScoreHomeExtraTime,
		match.ScoreAwayExtraTime,
		match.HasPenalties,
		match.ScoreHomePenalties,
		match.ScoreAwayPenalties,
		match.WinnerTeamID,
		match.Phase,
		match.ParentMatchID,
		match.LeftChildMatchID,
		match.RightChildMatchID,
		time.Now(),
		match.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were updated")
	}

	return nil
}
