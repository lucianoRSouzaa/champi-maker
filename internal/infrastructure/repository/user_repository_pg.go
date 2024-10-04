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

type userRepositoryPg struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPg(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepositoryPg{pool: pool}
}

func (r *userRepositoryPg) Create(ctx context.Context, user *entity.User) error {
	query := `
        INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *userRepositoryPg) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
        SELECT id, name, email, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Usuário não encontrado
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryPg) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
        SELECT id, name, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	row := r.pool.QueryRow(ctx, query, email)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Usuário não encontrado
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryPg) Update(ctx context.Context, user *entity.User) error {
	query := `
        UPDATE users
        SET name = $1,
            email = $2,
            password_hash = $3,
            updated_at = $4
        WHERE id = $5
    `
	commandTag, err := r.pool.Exec(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *userRepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM users
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
