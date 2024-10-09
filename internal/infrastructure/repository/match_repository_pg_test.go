package repository

import (
	"champi-maker/internal/domain/entity"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createUser(userID uuid.UUID, pool *pgxpool.Pool) (uuid.UUID, error) {
	ctx := context.Background()

	// criar user
	user := &entity.User{
		ID:           userID,
		Name:         "Usuário Teste",
		Email:        "teste",
		PasswordHash: "teste",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		INSERT INTO users (
			id, name, email, password_hash, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6 
		)
	`

	_, err := pool.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

	return userID, err
}

func createChampionship(championshipID uuid.UUID, pool *pgxpool.Pool) (uuid.UUID, error) {
	ctx := context.Background()

	championship := &entity.Championship{
		ID:               championshipID,
		Name:             "Campeonato Teste",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		Phases:           3,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `
		INSERT INTO championships (
			id, name, type, tiebreaker_method, phases, progression_type, created_at, updated_at
		) VALUES (
		 	$1, $2, $3, $4, $5, $6, $7, $8						
		)
	`

	_, err := pool.Exec(ctx, query, championship.ID, championship.Name, championship.Type, championship.TiebreakerMethod, championship.Phases, championship.ProgressionType, championship.CreatedAt, championship.UpdatedAt)

	return championshipID, err
}

func createTeam(teamID, userID uuid.UUID, pool *pgxpool.Pool) (uuid.UUID, error) {
	ctx := context.Background()

	team := &entity.Team{
		ID:        teamID,
		Name:      "Corinthians",
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO teams (
			id, name, user_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5 
		)
	`

	_, err := pool.Exec(ctx, query, team.ID, team.Name, team.UserID, team.CreatedAt, team.UpdatedAt)

	return teamID, err
}

func TestMatchRepositoryPg_CreateAndGetByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	homeTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	awayTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	now := time.Now()

	match := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &now,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Teste de criação
	err = matchRepo.Create(ctx, match)
	require.NoError(t, err)

	// Teste de recuperação pelo ID
	retrievedMatch, err := matchRepo.GetByID(ctx, match.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedMatch)

	// Verificações
	assert.Equal(t, match.ID, retrievedMatch.ID)
	assert.Equal(t, match.ChampionshipID, retrievedMatch.ChampionshipID)
	assert.Equal(t, match.HomeTeamID, retrievedMatch.HomeTeamID)
	assert.Equal(t, match.AwayTeamID, retrievedMatch.AwayTeamID)
	assert.Equal(t, match.Status, retrievedMatch.Status)
}

func TestMatchRepositoryPg_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	homeTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	awayTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	now := time.Now()

	match := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &now,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Teste de criação
	err = matchRepo.Create(ctx, match)
	require.NoError(t, err)

	// Atualizar o match
	match.Status = entity.MatchStatusFinished
	match.ScoreHome = 2
	match.ScoreAway = 1
	match.WinnerTeamID = match.HomeTeamID
	match.UpdatedAt = time.Now()

	err = matchRepo.Update(ctx, match)
	require.NoError(t, err)

	// Recuperar o match atualizado
	updatedMatch, err := matchRepo.GetByID(ctx, match.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedMatch)

	// Verificações
	assert.Equal(t, entity.MatchStatusFinished, updatedMatch.Status)
	assert.Equal(t, 2, updatedMatch.ScoreHome)
	assert.Equal(t, 1, updatedMatch.ScoreAway)
	assert.Equal(t, match.WinnerTeamID, updatedMatch.WinnerTeamID)
}

func TestMatchRepositoryPg_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	homeTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	awayTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	now := time.Now()

	match := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &now,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = matchRepo.Create(ctx, match)
	require.NoError(t, err)

	// Deletar o match
	err = matchRepo.Delete(ctx, match.ID)
	require.NoError(t, err)

	// Tentar recuperar o match deletado
	deletedMatch, err := matchRepo.GetByID(ctx, match.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedMatch)
}

func TestMatchRepositoryPg_GetByChampionshipID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	homeTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	awayTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	// Criar vários matches para o mesmo campeonato
	matchDate1 := time.Now().Add(-2 * time.Hour)
	match1 := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &matchDate1,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	matchDate2 := time.Now().Add(-1 * time.Hour)
	match2 := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &matchDate2,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = matchRepo.Create(ctx, match1)
	require.NoError(t, err)

	err = matchRepo.Create(ctx, match2)
	require.NoError(t, err)

	// Recuperar os matches pelo ChampionshipID
	matches, err := matchRepo.GetByChampionshipID(ctx, championshipId)
	require.NoError(t, err)
	require.Len(t, matches, 2)

	// Verificar se os matches são os esperados
	assert.Equal(t, match1.ID, matches[0].ID)
	assert.Equal(t, match2.ID, matches[1].ID)
}

func TestMatchRepositoryPg_GetByPhase(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	championshipId, err := createChampionship(uuid.New(), pool)
	require.NoError(t, err)

	homeTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	awayTeamID, err := createTeam(uuid.New(), userID, pool)
	require.NoError(t, err)

	// Criar vários matches em fases diferentes
	matchDate1 := time.Now().Add(-2 * time.Hour)
	match1 := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &matchDate1,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	matchDate2 := time.Now().Add(-1 * time.Hour)
	match2 := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championshipId,
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &matchDate2,
		Status:         entity.MatchStatusScheduled,
		Phase:          2,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = matchRepo.Create(ctx, match1)
	require.NoError(t, err)

	err = matchRepo.Create(ctx, match2)
	require.NoError(t, err)

	// Recuperar os matches da fase 1
	matchesPhase1, err := matchRepo.GetByPhase(ctx, championshipId, 1)
	require.NoError(t, err)
	require.Len(t, matchesPhase1, 1)
	assert.Equal(t, match1.ID, matchesPhase1[0].ID)

	// Recuperar os matches da fase 2
	matchesPhase2, err := matchRepo.GetByPhase(ctx, championshipId, 2)
	require.NoError(t, err)
	require.Len(t, matchesPhase2, 1)
	assert.Equal(t, match2.ID, matchesPhase2[0].ID)
}

func TestMatchRepositoryPg_Update_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	// Criar um match que não existe no banco
	homeTeamID := uuid.New()
	awayTeamID := uuid.New()
	matchDate := time.Now()

	match := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &matchDate,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		UpdatedAt:      time.Now(),
	}

	// Tentar atualizar
	err := matchRepo.Update(ctx, match)
	require.Error(t, err)
	assert.Equal(t, "no rows were updated", err.Error())
}

func TestMatchRepositoryPg_Delete_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	matchRepo := NewMatchRepositoryPg(pool)

	// ID aleatório que não existe
	nonExistentID := uuid.New()

	// Tentar deletar
	err := matchRepo.Delete(ctx, nonExistentID)
	require.Error(t, err)
	assert.Equal(t, "no rows were deleted", err.Error())
}
