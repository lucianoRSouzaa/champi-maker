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
	"golang.org/x/crypto/bcrypt"
)

type MockTokenProvider struct{}

func (m *MockTokenProvider) GenerateToken(userID uuid.UUID) (string, error) {
	return "mocked-jwt-token", nil
}

func hashPassword(password string) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashed := string(bytes)
	return &hashed, nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func stringPtr(s string) *string {
	return &s
}

func TestUserHandler_Register_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/register", userHandler.Register)

	requestBody := handler.RegisterUserRequest{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response entity.User
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, requestBody.Name, response.Name)
	assert.Equal(t, requestBody.Email, response.Email)
	assert.NotEqual(t, uuid.Nil, response.ID)
	assert.Equal(t, "", response.PasswordHash)
}

func TestUserHandler_Register_MissingFields(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/register", userHandler.Register)

	requestBody := map[string]interface{}{
		"name": "Test User",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Key: 'RegisterUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag")
	assert.Contains(t, recorder.Body.String(), "Key: 'RegisterUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag")
}

func TestUserHandler_Register_EmailAlreadyInUse(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	existingUser := &entity.User{
		ID:           uuid.New(),
		Name:         "Existing User",
		Email:        "existing@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, existingUser)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/register", userHandler.Register)

	requestBody := handler.RegisterUserRequest{
		Name:     "New User",
		Email:    "existing@example.com",
		Password: "newpassword",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "email já está em uso")
}

func TestUserHandler_Login_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	password := "password123"
	hashedPassword, err := hashPassword(password)
	require.NoError(t, err)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Login User",
		Email:        "loginuser@example.com",
		PasswordHash: *hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/login", userHandler.Login)

	requestBody := handler.LoginUserRequest{
		Email:    "loginuser@example.com",
		Password: "password123",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	token, exists := response["token"]
	assert.True(t, exists)
	assert.Equal(t, "mocked-jwt-token", token)
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	password := "password123"
	hashedPassword, err := hashPassword(password)
	require.NoError(t, err)

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Login User",
		Email:        "loginuser@example.com",
		PasswordHash: *hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/login", userHandler.Login)

	requestBody := handler.LoginUserRequest{
		Email:    "loginuser@example.com",
		Password: "wrongpassword",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "credenciais inválidas")
}

func TestUserHandler_Login_MissingFields(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/login", userHandler.Login)

	requestBody := map[string]interface{}{
		"password": "password123",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Key: 'LoginUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag")
}

func TestUserHandler_GetProfile_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Profile User",
		Email:        "profileuser@example.com",
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

	router.GET("/users/profile", testAuthMiddleware, userHandler.GetProfile)

	req, err := http.NewRequest(http.MethodGet, "/users/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer mocked-jwt-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response entity.User
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, "", response.PasswordHash) // PasswordHash não deve ser retornado
}

func TestUserHandler_GetProfile_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userHandler := handler.NewUserHandler(nil) // Serviço não é necessário neste teste

	router := gin.Default()
	router.GET("/users/profile", userHandler.GetProfile)

	req, err := http.NewRequest(http.MethodGet, "/users/profile", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "ID de usuário não encontrado no contexto")
}

func TestUserHandler_GetProfile_UserNotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		nonExistentID := uuid.New()
		c.Set("userID", nonExistentID)
		c.Next()
	}

	router.GET("/users/profile", testAuthMiddleware, userHandler.GetProfile)

	req, err := http.NewRequest(http.MethodGet, "/users/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer mocked-jwt-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "usuário não encontrado")
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Old Name",
		Email:        "oldemail@example.com",
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

	router.PUT("/users/profile", testAuthMiddleware, userHandler.UpdateProfile)

	updateData := handler.UpdateUserRequest{
		Name:     stringPtr("New Name"),
		Email:    stringPtr("newemail@example.com"),
		Password: stringPtr("newpassword123"),
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/users/profile", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Perfil atualizado com sucesso")

	updatedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "newemail@example.com", updatedUser.Email)
}

func TestUserHandler_UpdateProfile_InvalidData(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Old Name",
		Email:        "oldemail@example.com",
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

	router.PUT("/users/profile", testAuthMiddleware, userHandler.UpdateProfile)

	updateData := handler.UpdateUserRequest{
		Email: stringPtr("invalid-email"),
	}

	jsonBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/users/profile", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Perfil atualizado com sucesso")

	updatedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "invalid-email", updatedUser.Email)
}

func TestUserHandler_DeleteAccount_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	ctx := context.Background()

	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Delete User",
		Email:        "deleteuser@example.com",
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

	router.DELETE("/users/profile", testAuthMiddleware, userHandler.DeleteAccount)

	req, err := http.NewRequest(http.MethodDelete, "/users/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer mocked-jwt-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Conta excluída com sucesso")

	deletedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserHandler_DeleteAccount_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	defer teardownTestDB(t, pool)

	userRepo := repository.NewUserRepositoryPg(pool)
	tokenProvider := &MockTokenProvider{}
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testAuthMiddleware := func(c *gin.Context) {
		nonExistentID := uuid.New()
		c.Set("userID", nonExistentID)
		c.Next()
	}

	router.DELETE("/users/profile", testAuthMiddleware, userHandler.DeleteAccount)

	req, err := http.NewRequest(http.MethodDelete, "/users/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer mocked-jwt-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "usuário não encontrado")
}
