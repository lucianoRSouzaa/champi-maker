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

type teamRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewTeamRepositoryPg(pool *pgxpool.Pool) repository.TeamRepository {
	return &teamRepositoryPg{pool: pool}
}

func (r *teamRepositoryPg) Create(ctx context.Context, team *entity.Team) error {
	query := `
        INSERT INTO teams (id, name, logo, user_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.pool.Exec(ctx, query,
		team.ID,
		team.Name,
		team.Logo,
		team.UserID,
		team.CreatedAt,
		team.UpdatedAt,
	)
	return err
}

func (r *teamRepositoryPg) GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	query := `
        SELECT id, name, logo, user_id, created_at, updated_at
        FROM teams
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)

	var team entity.Team
	err := row.Scan(
		&team.ID,
		&team.Name,
		&team.Logo,
		&team.UserID,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Time não encontrado
		}
		return nil, err
	}

	return &team, nil
}

func (r *teamRepositoryPg) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Team, error) {
	query := `
        SELECT id, name, logo, user_id, created_at, updated_at
        FROM teams
        WHERE user_id = $1
    `
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*entity.Team
	for rows.Next() {
		var team entity.Team
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Logo,
			&team.UserID,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}

	return teams, nil
}

func (r *teamRepositoryPg) Update(ctx context.Context, team *entity.Team) error {
	query := `
        UPDATE teams
        SET name = $1,
            logo = $2,
            updated_at = $3
        WHERE id = $4
    `
	commandTag, err := r.pool.Exec(ctx, query,
		team.Name,
		team.Logo,
		time.Now(),
		team.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *teamRepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM teams
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

func (r *teamRepositoryPg) List(ctx context.Context) ([]*entity.Team, error) {
	query := `
        SELECT id, name, logo, user_id, created_at, updated_at
        FROM teams
    `
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*entity.Team
	for rows.Next() {
		var team entity.Team
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Logo,
			&team.UserID,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}

	return teams, nil
}
