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

func TestTeamRepositoryPg_CreateAndGetByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	// Criar um time de teste
	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time Teste",
		Logo:      "https://example.com/logo.png",
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Teste de criação
	err = teamRepo.Create(ctx, team)
	require.NoError(t, err)

	// Teste de recuperação pelo ID
	retrievedTeam, err := teamRepo.GetByID(ctx, team.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTeam)

	// Verificações
	assert.Equal(t, team.ID, retrievedTeam.ID)
	assert.Equal(t, team.Name, retrievedTeam.Name)
	assert.Equal(t, team.Logo, retrievedTeam.Logo)
	assert.Equal(t, team.UserID, retrievedTeam.UserID)
}

func createUserWithEmail(userID uuid.UUID, userEmail string, pool *pgxpool.Pool) (uuid.UUID, error) {
	ctx := context.Background()

	// criar user
	user := &entity.User{
		ID:           userID,
		Name:         "Usuário Teste",
		Email:        userEmail,
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

func TestTeamRepositoryPg_GetByUserID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	// Criar dois usuários de teste
	user1, err := createUserWithEmail(uuid.New(), "teste 1", pool)
	require.NoError(t, err)

	user2, err := createUserWithEmail(uuid.New(), "teste 2", pool)
	require.NoError(t, err)

	// Criar times para cada usuário
	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time User1",
		Logo:      "https://example.com/logo1.png",
		UserID:    user1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time User2",
		Logo:      "https://example.com/logo2.png",
		UserID:    user2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	team3 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time User1 (2)",
		Logo:      "https://example.com/logo3.png",
		UserID:    user1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Inserir os times
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)

	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	err = teamRepo.Create(ctx, team3)
	require.NoError(t, err)

	// Recuperar times do user1
	teamsUser1, err := teamRepo.GetByUserID(ctx, user1)
	require.NoError(t, err)
	require.Len(t, teamsUser1, 2)
	assert.Equal(t, team1.ID, teamsUser1[0].ID)
	assert.Equal(t, team3.ID, teamsUser1[1].ID)

	// Recuperar times do user2
	teamsUser2, err := teamRepo.GetByUserID(ctx, user2)
	require.NoError(t, err)
	require.Len(t, teamsUser2, 1)
	assert.Equal(t, team2.ID, teamsUser2[0].ID)
}

func TestTeamRepositoryPg_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	// Criar um time de teste
	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time Original",
		Logo:      "https://example.com/logo.png",
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = teamRepo.Create(ctx, team)
	require.NoError(t, err)

	// Atualizar o time
	team.Name = "Time Atualizado"
	team.Logo = "https://example.com/new_logo.png"
	team.UpdatedAt = time.Now()

	err = teamRepo.Update(ctx, team)
	require.NoError(t, err)

	// Recuperar o time atualizado
	updatedTeam, err := teamRepo.GetByID(ctx, team.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedTeam)

	// Verificações
	assert.Equal(t, "Time Atualizado", updatedTeam.Name)
	assert.Equal(t, "https://example.com/new_logo.png", updatedTeam.Logo)
}

func TestTeamRepositoryPg_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	// Criar um time de teste
	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time para Deletar",
		Logo:      "https://example.com/logo.png",
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = teamRepo.Create(ctx, team)
	require.NoError(t, err)

	// Deletar o time
	err = teamRepo.Delete(ctx, team.ID)
	require.NoError(t, err)

	// Tentar recuperar o time deletado
	deletedTeam, err := teamRepo.GetByID(ctx, team.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedTeam)
}

func TestTeamRepositoryPg_List(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	userID, err := createUser(uuid.New(), pool)
	require.NoError(t, err)

	// Criar vários times
	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time 1",
		Logo:      "https://example.com/logo1.png",
		UserID:    userID,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time 2",
		Logo:      "https://example.com/logo2.png",
		UserID:    userID,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	// Listar os times
	teams, err := teamRepo.List(ctx)
	require.NoError(t, err)
	require.Len(t, teams, 2)

	// Verificar se os times estão na lista
	teamIDs := []uuid.UUID{teams[0].ID, teams[1].ID}
	assert.Contains(t, teamIDs, team1.ID)
	assert.Contains(t, teamIDs, team2.ID)
}

func TestTeamRepositoryPg_Update_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	// Criar um time que não existe no banco
	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Time Inexistente",
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(), // ID de usuário aleatório
		UpdatedAt: time.Now(),
	}

	// Tentar atualizar
	err := teamRepo.Update(ctx, team)
	require.Error(t, err)
	assert.Equal(t, "no rows were updated", err.Error())
}

func TestTeamRepositoryPg_Delete_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	teamRepo := NewTeamRepositoryPg(pool)

	// ID aleatório que não existe
	nonExistentID := uuid.New()

	// Tentar deletar
	err := teamRepo.Delete(ctx, nonExistentID)
	require.Error(t, err)
	assert.Equal(t, "no rows were deleted", err.Error())
}
