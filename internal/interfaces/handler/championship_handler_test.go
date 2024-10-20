package handler_test

import (
	"bytes"
	"champi-maker/internal/application/service"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/infrastructure/config"
	"champi-maker/internal/infrastructure/repository"
	"champi-maker/internal/interfaces/handler"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

type MockMessagePublisher struct{}

func (m *MockMessagePublisher) PublishChampionshipCreated(ctx context.Context, championshipID uuid.UUID, teamIDs []uuid.UUID) error {
	return nil
}

func TestChampionshipHandler_CreateChampionship_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)

	ctx := context.Background()

	userID := uuid.New()
	user := &entity.User{
		ID:           userID,
		Name:         "User",
		Email:        "teste@gmail.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo := repository.NewUserRepositoryPg(pool)

	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	championshipHandler := handler.NewChampionshipHandler(championshipService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/championships", championshipHandler.CreateChampionship)

	requestBody := handler.CreateChampionshipRequest{
		Name:             "Test Championship",
		Type:             "league",
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
		TeamIDs:          []uuid.UUID{team1.ID, team2.ID},
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/championships", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response entity.Championship
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, requestBody.Name, response.Name)
	assert.Equal(t, entity.ChampionshipType(requestBody.Type), response.Type)
	assert.Equal(t, entity.TiebreakerMethod(requestBody.TiebreakerMethod), response.TiebreakerMethod)
	assert.Equal(t, entity.ProgressionType(requestBody.ProgressionType), response.ProgressionType)
	assert.NotEqual(t, uuid.Nil, response.ID)
}

func TestChampionshipHandler_CreateChampionship_InvalidData(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/championships", championshipHandler.CreateChampionship)

	// Requisição com dados inválidos (faltando o campo "name")
	requestBody := map[string]interface{}{
		"type":              "league",
		"tiebreaker_method": "points",
		"progression_type":  "fixed",
		"team_ids":          []string{},
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/championships", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Key: 'CreateChampionshipRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag")
}

func TestChampionshipHandler_CreateChampionship_TeamNotFound(t *testing.T) {
	// Configuração do banco de dados
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/championships", championshipHandler.CreateChampionship)

	// IDs de times que não existem
	teamID1 := uuid.New()
	teamID2 := uuid.New()

	requestBody := handler.CreateChampionshipRequest{
		Name:             "Test Championship",
		Type:             "league",
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
		TeamIDs:          []uuid.UUID{teamID1, teamID2},
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/championships", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "team with ID")
}

func TestChampionshipHandler_GetChampionshipByID_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test Championship",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championshipRepo.Create(ctx, championship)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/championships/:id", championshipHandler.GetChampionshipByID)

	req, err := http.NewRequest(http.MethodGet, "/championships/"+championship.ID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response entity.Championship
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, championship.ID, response.ID)
	assert.Equal(t, championship.Name, response.Name)
}

func TestChampionshipHandler_GetChampionshipByID_InvalidID(t *testing.T) {
	championshipHandler := handler.NewChampionshipHandler(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/championships/:id", championshipHandler.GetChampionshipByID)

	req, err := http.NewRequest(http.MethodGet, "/championships/invalid-uuid", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "invalid championship ID")
}

func TestChampionshipHandler_GetChampionshipByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/championships/:id", championshipHandler.GetChampionshipByID)

	nonExistentID := uuid.New()

	req, err := http.NewRequest(http.MethodGet, "/championships/"+nonExistentID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "championship not found")
}

func TestChampionshipHandler_UpdateChampionship_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Original Name",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championshipRepo.Create(ctx, championship)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/championships/:id", championshipHandler.UpdateChampionship)

	// Dados de atualização
	updateData := handler.UpdateChampionshipRequest{
		Name:             "Updated Name",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/championships/"+championship.ID.String(), bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// fazer get by id repository para verificar se foi atualizado
	updatedChampionship, err := championshipRepo.GetByID(ctx, championship.ID)

	assert.Equal(t, updateData.Name, updatedChampionship.Name)
	assert.Equal(t, updateData.Type, updatedChampionship.Type)
	assert.Equal(t, updateData.TiebreakerMethod, updatedChampionship.TiebreakerMethod)
	assert.Equal(t, updateData.ProgressionType, updatedChampionship.ProgressionType)
}

func TestChampionshipHandler_UpdateChampionship_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/championships/:id", championshipHandler.UpdateChampionship)

	updateData := handler.CreateChampionshipRequest{
		Name:             "Updated Name",
		Type:             "league",
		TiebreakerMethod: "points",
		ProgressionType:  "fixed",
		TeamIDs:          []uuid.UUID{},
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	nonExistentID := uuid.New()

	req, err := http.NewRequest(http.MethodPut, "/championships/"+nonExistentID.String(), bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "championship with ID")
}

func TestChampionshipHandler_DeleteChampionship_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test Championship",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championshipRepo.Create(ctx, championship)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/championships/:id", championshipHandler.DeleteChampionship)

	req, err := http.NewRequest(http.MethodDelete, "/championships/"+championship.ID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "championship deleted successfully")
}

func TestChampionshipHandler_DeleteChampionship_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/championships/:id", championshipHandler.DeleteChampionship)

	nonExistentID := uuid.New()

	req, err := http.NewRequest(http.MethodDelete, "/championships/"+nonExistentID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "championship with ID")
}

func TestChampionshipHandler_ListChampionships_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	messagePublisher := &MockMessagePublisher{}

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	// Criar campeonatos de teste
	championship1 := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Championship 1",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	championship2 := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Championship 2",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionRandomDraw,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championshipRepo.Create(ctx, championship1)
	require.NoError(t, err)
	err = championshipRepo.Create(ctx, championship2)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/championships", championshipHandler.ListChampionships)

	req, err := http.NewRequest(http.MethodGet, "/championships", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []entity.Championship
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	names := []string{response[0].Name, response[1].Name}
	assert.Contains(t, names, "Championship 1")
	assert.Contains(t, names, "Championship 2")
}
