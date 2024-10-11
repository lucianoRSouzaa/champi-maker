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

func TestUserRepositoryPg_CreateAndGetByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	userRepo := NewUserRepositoryPg(pool)

	// Criar um usuário de teste
	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Teste de criação
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Teste de recuperação pelo ID
	retrievedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	// Verificações
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.PasswordHash, retrievedUser.PasswordHash)
}

func TestUserRepositoryPg_GetByEmail(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	userRepo := NewUserRepositoryPg(pool)

	// Criar um usuário de teste
	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Teste de recuperação pelo Email
	retrievedUser, err := userRepo.GetByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	// Verificações
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.PasswordHash, retrievedUser.PasswordHash)
}

func TestUserRepositoryPg_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	userRepo := NewUserRepositoryPg(pool)

	// Criar um usuário de teste
	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Original User",
		Email:        "original@example.com",
		PasswordHash: "originalhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Atualizar o usuário
	user.Name = "Updated User"
	user.Email = "updated@example.com"
	user.PasswordHash = "updatedhash"
	user.UpdatedAt = time.Now()

	err = userRepo.Update(ctx, user)
	require.NoError(t, err)

	// Recuperar o usuário atualizado
	updatedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)

	// Verificações
	assert.Equal(t, "Updated User", updatedUser.Name)
	assert.Equal(t, "updated@example.com", updatedUser.Email)
	assert.Equal(t, "updatedhash", updatedUser.PasswordHash)
}

func TestUserRepositoryPg_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	userRepo := NewUserRepositoryPg(pool)

	// Criar um usuário de teste
	user := &entity.User{
		ID:           uuid.New(),
		Name:         "User to Delete",
		Email:        "delete@example.com",
		PasswordHash: "deletehash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Deletar o usuário
	err = userRepo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// Tentar recuperar o usuário deletado
	deletedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserRepositoryPg_Update_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	userRepo := NewUserRepositoryPg(pool)

	// Criar um usuário que não existe no banco
	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Nonexistent User",
		Email:        "nonexistent@example.com",
		PasswordHash: "hash",
		UpdatedAt:    time.Now(),
	}

	// Tentar atualizar
	err := userRepo.Update(ctx, user)
	require.Error(t, err)
	assert.Equal(t, "no rows were updated", err.Error())
}

func TestUserRepositoryPg_Delete_NonExistent(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	userRepo := NewUserRepositoryPg(pool)

	// ID aleatório que não existe
	nonExistentID := uuid.New()

	// Tentar deletar
	err := userRepo.Delete(ctx, nonExistentID)
	require.Error(t, err)
	assert.Equal(t, "no rows were deleted", err.Error())
}
