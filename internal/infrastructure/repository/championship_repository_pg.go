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

type championshipRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewChampionshipRepositoryPg(pool *pgxpool.Pool) repository.ChampionshipRepository {
	return &championshipRepositoryPg{pool: pool}
}

func (r *championshipRepositoryPg) Create(ctx context.Context, championship *entity.Championship) error {
	query := `
        INSERT INTO championships (id, name, type, tiebreaker_method, progression_type, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.pool.Exec(ctx, query,
		championship.ID,
		championship.Name,
		string(championship.Type),
		string(championship.TiebreakerMethod),
		string(championship.ProgressionType),
		championship.CreatedAt,
		championship.UpdatedAt,
	)
	return err
}

func (r *championshipRepositoryPg) GetByID(ctx context.Context, id uuid.UUID) (*entity.Championship, error) {
	query := `
        SELECT id, name, type, tiebreaker_method, progression_type, created_at, updated_at
        FROM championships
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)

	var championship entity.Championship
	var typeStr, tiebreakerMethodStr, progressionTypeStr string

	err := row.Scan(
		&championship.ID,
		&championship.Name,
		&typeStr,
		&tiebreakerMethodStr,
		&progressionTypeStr,
		&championship.CreatedAt,
		&championship.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Campeonato não encontrado
		}
		return nil, err
	}

	// Converter strings para tipos enumerados
	championship.Type = entity.ChampionshipType(typeStr)
	championship.TiebreakerMethod = entity.TiebreakerMethod(tiebreakerMethodStr)
	championship.ProgressionType = entity.ProgressionType(progressionTypeStr)

	return &championship, nil
}

func (r *championshipRepositoryPg) Update(ctx context.Context, championship *entity.Championship) error {
	query := `
        UPDATE championships
        SET name = $1,
            type = $2,
            tiebreaker_method = $3,
            progression_type = $4,
            updated_at = $5
        WHERE id = $6
    `
	commandTag, err := r.pool.Exec(ctx, query,
		championship.Name,
		string(championship.Type),
		string(championship.TiebreakerMethod),
		string(championship.ProgressionType),
		time.Now(),
		championship.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *championshipRepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM championships
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

func (r *championshipRepositoryPg) List(ctx context.Context) ([]*entity.Championship, error) {
	query := `
        SELECT id, name, type, tiebreaker_method, progression_type, created_at, updated_at
        FROM championships
        ORDER BY created_at DESC
    `
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var championships []*entity.Championship

	for rows.Next() {
		var championship entity.Championship
		var typeStr, tiebreakerMethodStr, progressionTypeStr string

		err := rows.Scan(
			&championship.ID,
			&championship.Name,
			&typeStr,
			&tiebreakerMethodStr,
			&progressionTypeStr,
			&championship.CreatedAt,
			&championship.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Converter strings para tipos enumerados
		championship.Type = entity.ChampionshipType(typeStr)
		championship.TiebreakerMethod = entity.TiebreakerMethod(tiebreakerMethodStr)
		championship.ProgressionType = entity.ProgressionType(progressionTypeStr)

		championships = append(championships, &championship)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return championships, nil
}
