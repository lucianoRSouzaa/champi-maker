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

func TestTeamHandler_CreateTeam_Success(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Team Owner",
		Email:        "teamowner@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Next()
	}

	router.POST("/teams", testAuthMiddleware, teamHandler.CreateTeam)

	requestBody := handler.CreateTeamRequest{
		Name: "New Team",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response entity.Team
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, requestBody.Name, response.Name)
	assert.Equal(t, user.ID, response.UserID)
	assert.NotEqual(t, uuid.Nil, response.ID)
}

func TestTeamHandler_CreateTeam_MissingFields(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", uuid.New())
		c.Next()
	}

	router.POST("/teams", testAuthMiddleware, teamHandler.CreateTeam)

	requestBody := map[string]interface{}{
		// "name" está ausente
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Key: 'CreateTeamRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag")
}

func TestTeamHandler_CreateTeam_UserNotFound(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", uuid.New())
		c.Next()
	}

	router.POST("/teams", testAuthMiddleware, teamHandler.CreateTeam)

	requestBody := handler.CreateTeamRequest{
		Name: "Team Without Owner",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "user with ID")
}

func TestTeamHandler_CreateTeam_NameAlreadyExists(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Team Owner",
		Email:        "teamowner@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	existingTeam := &entity.Team{
		ID:        uuid.New(),
		Name:      "Existing Team",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = teamRepo.Create(ctx, existingTeam)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Next()
	}

	router.POST("/teams", testAuthMiddleware, teamHandler.CreateTeam)

	requestBody := handler.CreateTeamRequest{
		Name: "Existing Team",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "team with name Existing Team already exists")
}

func TestTeamHandler_GetTeamByID_Success(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Team Owner",
		Email:        "teamowner@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Unique Team",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = teamRepo.Create(ctx, team)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/teams/:id", teamHandler.GetTeamByID)

	req, err := http.NewRequest(http.MethodGet, "/teams/"+team.ID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response entity.Team
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, team.ID, response.ID)
	assert.Equal(t, team.Name, response.Name)
	assert.Equal(t, team.UserID, response.UserID)
}

func TestTeamHandler_GetTeamByID_InvalidID(t *testing.T) {

	gin.SetMode(gin.TestMode)
	teamHandler := handler.NewTeamHandler(nil)

	router := gin.Default()
	router.GET("/teams/:id", teamHandler.GetTeamByID)

	req, err := http.NewRequest(http.MethodGet, "/teams/invalid-uuid", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "ID de time inválido")
}

func TestTeamHandler_GetTeamByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/teams/:id", teamHandler.GetTeamByID)

	nonExistentID := uuid.New()
	req, err := http.NewRequest(http.MethodGet, "/teams/"+nonExistentID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "time não encontrado")
}

func TestTeamHandler_UpdateTeam_Success(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Team Owner",
		Email:        "teamowner@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Original Team",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = teamRepo.Create(ctx, team)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/teams/:id", teamHandler.UpdateTeam)

	updateData := handler.UpdateTeamRequest{
		Name: "Updated Team Name",
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/teams/"+team.ID.String(), bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response entity.Team
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, team.ID, response.ID)
	assert.Equal(t, "Updated Team Name", response.Name)

	updatedTeam, err := teamRepo.GetByID(ctx, team.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Team Name", updatedTeam.Name)
}

func TestTeamHandler_UpdateTeam_InvalidID(t *testing.T) {

	gin.SetMode(gin.TestMode)
	teamHandler := handler.NewTeamHandler(nil)

	router := gin.Default()
	router.PUT("/teams/:id", teamHandler.UpdateTeam)

	requestBody := handler.UpdateTeamRequest{
		Name: "New Name",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/teams/invalid-uuid", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "ID de time inválido")
}

func TestTeamHandler_UpdateTeam_NotFound(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/teams/:id", teamHandler.UpdateTeam)

	nonExistentID := uuid.New()
	updateData := handler.UpdateTeamRequest{
		Name: "Non-Existent Team",
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/teams/"+nonExistentID.String(), bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "team with ID")
}

func TestTeamHandler_DeleteTeam_Success(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Team Owner",
		Email:        "teamowner@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	team := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team To Delete",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = teamRepo.Create(ctx, team)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/teams/:id", teamHandler.DeleteTeam)

	req, err := http.NewRequest(http.MethodDelete, "/teams/"+team.ID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "time excluído com sucesso")

	deletedTeam, err := teamRepo.GetByID(ctx, team.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedTeam)
}

func TestTeamHandler_DeleteTeam_NotFound(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/teams/:id", teamHandler.DeleteTeam)

	nonExistentID := uuid.New()
	req, err := http.NewRequest(http.MethodDelete, "/teams/"+nonExistentID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "team with ID")
}

func TestTeamHandler_ListTeams_Success(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Team Owner",
		Email:        "teamowner@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team One",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team Two",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/teams", teamHandler.ListTeams)

	req, err := http.NewRequest(http.MethodGet, "/teams", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []*entity.Team
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	names := []string{response[0].Name, response[1].Name}
	assert.Contains(t, names, "Team One")
	assert.Contains(t, names, "Team Two")
}

func TestTeamHandler_ListTeamsByUserID_Success(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	ctx := context.Background()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	user1 := &entity.User{
		ID:           uuid.New(),
		Name:         "User One",
		Email:        "userone@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	user2 := &entity.User{
		ID:           uuid.New(),
		Name:         "User Two",
		Email:        "usertwo@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user1)
	require.NoError(t, err)
	err = userRepo.Create(ctx, user2)
	require.NoError(t, err)

	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "User1 Team1",
		UserID:    user1.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "User1 Team2",
		UserID:    user1.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	team3 := &entity.Team{
		ID:        uuid.New(),
		Name:      "User2 Team1",
		UserID:    user2.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team3)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/teams/user/:user_id", teamHandler.ListTeamsByUserID)

	req, err := http.NewRequest(http.MethodGet, "/teams/user/"+user1.ID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []*entity.Team
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	names := []string{response[0].Name, response[1].Name}
	assert.Contains(t, names, "User1 Team1")
	assert.Contains(t, names, "User1 Team2")
}

func TestTeamHandler_ListTeamsByUserID_UserNotFound(t *testing.T) {

	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamHandler := handler.NewTeamHandler(teamService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/teams/user/:user_id", teamHandler.ListTeamsByUserID)

	nonExistentID := uuid.New()
	req, err := http.NewRequest(http.MethodGet, "/teams/user/"+nonExistentID.String(), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "user with ID")
}
