package repository

import (
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/infrastructure/config"

	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	err := config.LoadEnv()
	require.NoError(t, err)

	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	if dbURL == "" {
		t.Fatal("DATABASE_URL_TEST is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)

	return pool
}

func teardownTestDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE championships CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE teams CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE users CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE matches CASCADE")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "TRUNCATE TABLE statistics CASCADE")
	require.NoError(t, err)
}

func TestChampionshipRepositoryPg_CreateAndGetByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := NewChampionshipRepositoryPg(pool)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Teste",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		Phases:           3,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Teste de criação
	err := championshipRepo.Create(ctx, championship)
	require.NoError(t, err)

	// Teste de recuperação pelo ID
	retrievedChampionship, err := championshipRepo.GetByID(ctx, championship.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedChampionship)

	// Verificações
	assert.Equal(t, championship.ID, retrievedChampionship.ID)
	assert.Equal(t, championship.Name, retrievedChampionship.Name)
	assert.Equal(t, championship.Type, retrievedChampionship.Type)
	assert.Equal(t, championship.TiebreakerMethod, retrievedChampionship.TiebreakerMethod)
	assert.Equal(t, championship.ProgressionType, retrievedChampionship.ProgressionType)
}

func TestChampionshipRepositoryPg_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := NewChampionshipRepositoryPg(pool)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Original",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		Phases:           2,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championshipRepo.Create(ctx, championship)
	require.NoError(t, err)

	// Atualizar o campeonato
	championship.Name = "Campeonato Atualizado"
	championship.Type = entity.ChampionshipTypeLeague
	championship.TiebreakerMethod = entity.TiebreakerExtraTime
	championship.ProgressionType = entity.ProgressionRandomDraw
	championship.UpdatedAt = time.Now()

	err = championshipRepo.Update(ctx, championship)
	require.NoError(t, err)

	// Recuperar o campeonato atualizado
	updatedChampionship, err := championshipRepo.GetByID(ctx, championship.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedChampionship)

	// Verificações
	assert.Equal(t, "Campeonato Atualizado", updatedChampionship.Name)
	assert.Equal(t, entity.ChampionshipTypeLeague, updatedChampionship.Type)
	assert.Equal(t, entity.TiebreakerExtraTime, updatedChampionship.TiebreakerMethod)
	assert.Equal(t, entity.ProgressionRandomDraw, updatedChampionship.ProgressionType)
}

func TestChampionshipRepositoryPg_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := NewChampionshipRepositoryPg(pool)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Campeonato para Deletar",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championshipRepo.Create(ctx, championship)
	require.NoError(t, err)

	// Deletar o campeonato
	err = championshipRepo.Delete(ctx, championship.ID)
	require.NoError(t, err)

	// Tentar recuperar o campeonato deletado
	deletedChampionship, err := championshipRepo.GetByID(ctx, championship.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedChampionship)
}

func TestChampionshipRepositoryPg_List(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := NewChampionshipRepositoryPg(pool)

	// Inserir vários campeonatos
	championship1 := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Campeonato 1",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now().Add(-2 * time.Hour),
		UpdatedAt:        time.Now().Add(-2 * time.Hour),
	}

	championship2 := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Campeonato 2",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionRandomDraw,
		CreatedAt:        time.Now().Add(-1 * time.Hour),
		UpdatedAt:        time.Now().Add(-1 * time.Hour),
	}

	err := championshipRepo.Create(ctx, championship1)
	require.NoError(t, err)

	err = championshipRepo.Create(ctx, championship2)
	require.NoError(t, err)

	// Listar os campeonatos
	championships, err := championshipRepo.List(ctx)
	require.NoError(t, err)
	require.Len(t, championships, 2)

	// Verificar a ordem (de acordo com o ORDER BY created_at DESC)
	assert.Equal(t, championship2.ID, championships[0].ID)
	assert.Equal(t, championship1.ID, championships[1].ID)
}

func TestChampionshipRepositoryPg_Update_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := NewChampionshipRepositoryPg(pool)

	// Criar um campeonato que não existe no banco
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Inexistente",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		UpdatedAt:        time.Now(),
	}

	// Tentar atualizar
	err := championshipRepo.Update(ctx, championship)
	require.Error(t, err)
	assert.Equal(t, "no rows were updated", err.Error())
}

func TestChampionshipRepositoryPg_Delete_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := NewChampionshipRepositoryPg(pool)

	// ID aleatório que não existe
	nonExistentID := uuid.New()

	// Tentar deletar
	err := championshipRepo.Delete(ctx, nonExistentID)
	require.Error(t, err)
	assert.Equal(t, "no rows were deleted", err.Error())
}
