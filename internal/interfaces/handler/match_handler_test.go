package handler_test

import (
	"bytes"
	"champi-maker/internal/application/service"
	"champi-maker/internal/domain/entity"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchHandler_GetMatchByID_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	statisticsRepo := repository.NewStatisticsRepositoryPg(pool)
	statisticsService := service.NewStatisticsService(statisticsRepo, championshipRepo, teamRepo)

	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo, statisticsService)
	matchHandler := handler.NewMatchHandler(matchService)

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

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "User",
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
		UserID:    user.ID,
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	match := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championship.ID,
		HomeTeamID:     &team1.ID,
		AwayTeamID:     &team2.ID,
		MatchDate:      nil,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = matchRepo.Create(ctx, match)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/matches/:id", matchHandler.GetMatchByID)

	req, err := http.NewRequest(http.MethodGet, "/matches/"+match.ID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response entity.Match
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, match.ID, response.ID)
	assert.Equal(t, match.ChampionshipID, response.ChampionshipID)
}

func TestMatchHandler_ListMatchesByChampionship_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	matchRepo := repository.NewMatchRepositoryPg(pool)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)
	statisticsRepo := repository.NewStatisticsRepositoryPg(pool)
	statisticsService := service.NewStatisticsService(statisticsRepo, championshipRepo, teamRepo)

	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo, statisticsService)
	matchHandler := handler.NewMatchHandler(matchService)

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

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "User",
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
		UserID:    user.ID,
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	match1 := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championship.ID,
		HomeTeamID:     &team1.ID,
		AwayTeamID:     &team2.ID,
		MatchDate:      nil,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	match2 := &entity.Match{
		ID:             uuid.New(),
		ChampionshipID: championship.ID,
		HomeTeamID:     &team2.ID,
		AwayTeamID:     &team1.ID,
		MatchDate:      nil,
		Status:         entity.MatchStatusScheduled,
		Phase:          1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = matchRepo.Create(ctx, match1)
	require.NoError(t, err)
	err = matchRepo.Create(ctx, match2)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/championships/:championship_id/matches", matchHandler.ListMatchesByChampionship)

	req, err := http.NewRequest(http.MethodGet, "/championships/"+championship.ID.String()+"/matches", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []entity.Match
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
}

func TestMatchHandler_UpdateMatchResult_InvalidData(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	matchRepo := repository.NewMatchRepositoryPg(pool)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	statisticsRepo := repository.NewStatisticsRepositoryPg(pool)
	statisticsService := service.NewStatisticsService(statisticsRepo, championshipRepo, teamRepo)

	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo, statisticsService)
	matchHandler := handler.NewMatchHandler(matchService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/matches/:id/result", matchHandler.UpdateMatchResult)

	updateData := map[string]interface{}{
		"score_home": "invalid",
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/matches/"+uuid.New().String()+"/result", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
