package repository

import (
	"champi-maker/internal/domain/entity"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatisticsRepositoryPg_CreateAndGetByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	statsRepo := NewStatisticsRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	teamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	// Criar estatística
	stats := &entity.Statistics{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		TeamID:         teamID,
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Teste de criação
	err = statsRepo.Create(ctx, stats)
	require.NoError(t, err)

	// Teste de recuperação pelo ID
	retrievedStats, err := statsRepo.GetByID(ctx, stats.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedStats)

	// Verificações
	assert.Equal(t, stats.ID, retrievedStats.ID)
	assert.Equal(t, stats.ChampionshipID, retrievedStats.ChampionshipID)
	assert.Equal(t, stats.TeamID, retrievedStats.TeamID)
	assert.Equal(t, stats.Points, retrievedStats.Points)
}

func TestStatisticsRepositoryPg_GetByChampionshipAndTeam(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	statsRepo := NewStatisticsRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	teamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	// Criar estatística
	stats := &entity.Statistics{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		TeamID:         teamID,
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = statsRepo.Create(ctx, stats)
	require.NoError(t, err)

	// Recuperar estatística pelo ChampionshipID e TeamID
	retrievedStats, err := statsRepo.GetByChampionshipAndTeam(ctx, championshipId, teamID)
	require.NoError(t, err)
	require.NotNil(t, retrievedStats)

	// Verificações
	assert.Equal(t, stats.ID, retrievedStats.ID)
	assert.Equal(t, stats.Points, retrievedStats.Points)
}

func TestStatisticsRepositoryPg_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	statsRepo := NewStatisticsRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	teamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	// Criar estatística
	stats := &entity.Statistics{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		TeamID:         teamID,
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = statsRepo.Create(ctx, stats)
	require.NoError(t, err)

	// Atualizar estatística
	stats.MatchesPlayed += 1
	stats.Wins += 1
	stats.GoalsFor += 2
	stats.GoalDifference += 2
	stats.Points += 3
	stats.UpdatedAt = time.Now()

	err = statsRepo.Update(ctx, stats)
	require.NoError(t, err)

	// Recuperar estatística atualizada
	updatedStats, err := statsRepo.GetByID(ctx, stats.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedStats)

	// Verificações
	assert.Equal(t, stats.MatchesPlayed, updatedStats.MatchesPlayed)
	assert.Equal(t, stats.Wins, updatedStats.Wins)
	assert.Equal(t, stats.GoalsFor, updatedStats.GoalsFor)
	assert.Equal(t, stats.Points, updatedStats.Points)
}

func TestStatisticsRepositoryPg_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	statsRepo := NewStatisticsRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	teamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	// Criar estatística
	stats := &entity.Statistics{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		TeamID:         teamID,
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = statsRepo.Create(ctx, stats)
	require.NoError(t, err)

	// Deletar estatística
	err = statsRepo.Delete(ctx, stats.ID)
	require.NoError(t, err)

	// Tentar recuperar estatística deletada
	deletedStats, err := statsRepo.GetByID(ctx, stats.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedStats)
}

func TestStatisticsRepositoryPg_ListByChampionship(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	statsRepo := NewStatisticsRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	team1, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	team2, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	// Criar estatísticas
	stats1 := &entity.Statistics{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		TeamID:         team1,
		MatchesPlayed:  10,
		Wins:           7,
		Draws:          2,
		Losses:         1,
		GoalsFor:       20,
		GoalsAgainst:   10,
		GoalDifference: 10,
		Points:         23,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	stats2 := &entity.Statistics{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		TeamID:         team2,
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          3,
		Losses:         1,
		GoalsFor:       18,
		GoalsAgainst:   8,
		GoalDifference: 10,
		Points:         21,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = statsRepo.Create(ctx, stats1)
	require.NoError(t, err)

	err = statsRepo.Create(ctx, stats2)
	require.NoError(t, err)

	// Listar estatísticas por campeonato
	statsList, err := statsRepo.ListByChampionship(ctx, championshipId)
	require.NoError(t, err)
	require.Len(t, statsList, 2)

	// Verificar ordem (de acordo com o ORDER BY points DESC, goal_difference DESC, goals_for DESC)
	assert.Equal(t, stats1.ID, statsList[0].ID)
	assert.Equal(t, stats2.ID, statsList[1].ID)
}
