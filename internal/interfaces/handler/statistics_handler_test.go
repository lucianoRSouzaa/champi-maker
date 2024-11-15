package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"champi-maker/internal/application/service"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/infrastructure/repository"
	"champi-maker/internal/interfaces/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatisticsHandler_GenerateInitialStatistics_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	statisticsRepo := repository.NewStatisticsRepositoryPg(pool)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	statisticsService := service.NewStatisticsService(statisticsRepo, championshipRepo, teamRepo)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)

	ctx := context.Background()

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

	userId := uuid.New()

	user := &entity.User{
		ID:           userId,
		Name:         "Test User",
		Email:        "email@email.co",
		PasswordHash: "password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
	}
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", uuid.New())
		c.Next()
	}

	router.POST("/championships/:id/statistics/init", testAuthMiddleware, statisticsHandler.GenerateInitialStatistics)

	requestBody := handler.GenerateStatisticsRequest{
		TeamIDs: []uuid.UUID{team1.ID, team2.ID},
	}
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/championships/"+championship.ID.String()+"/statistics/init", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Estatísticas iniciais geradas com sucesso")

	statsList, err := statisticsRepo.ListByChampionship(ctx, championship.ID)
	require.NoError(t, err)
	assert.Len(t, statsList, 2)
}
