package service_test

import (
	"champi-maker/internal/application/service"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/infrastructure/config"
	"champi-maker/internal/infrastructure/messaging"
	"champi-maker/internal/infrastructure/repository"
	"fmt"

	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupDatabase(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE matches, championships, teams, users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestEndToEndChampionshipCreation(t *testing.T) {
	ctx := context.Background()

	err := config.LoadEnv()
	require.NoError(t, err)

	// Configurar o banco de dados de teste
	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	require.NotEmpty(t, dbURL, "DATABASE_URL_TEST is not set")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Limpar tabelas relevantes
	cleanupDatabase(t, pool)

	// Configurar o RabbitMQ de teste
	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL_TEST")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL_TEST is not set")

	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	require.NoError(t, err)
	defer rabbitConn.Close()

	// Inicializar repositórios
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	// Inicializar produtor de mensagens
	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created_test")
	require.NoError(t, err)

	// Inicializar serviços
	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo)

	// Iniciar consumidor
	rabbitMQConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created_test", matchService)
	require.NoError(t, err)

	// Iniciar consumo em uma goroutine
	go func() {
		err := rabbitMQConsumer.StartConsuming()
		require.NoError(t, err)
	}()

	userID := uuid.New()

	user := &entity.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test@gmail.com",
		PasswordHash: "123456",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.Create(ctx, user)

	// Criar times de teste
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
	team3 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team 3",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}

	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team3)
	require.NoError(t, err)

	// Criar um campeonato de teste
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test Championship",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	teamIDs := []uuid.UUID{team1.ID, team2.ID, team3.ID}

	err = championshipService.CreateChampionship(ctx, championship, teamIDs)
	require.NoError(t, err)

	// Aguardar um tempo para o consumidor processar a mensagem
	time.Sleep(5 * time.Second)

	// Verificar se as partidas foram geradas
	matches, err := matchRepo.GetByChampionshipID(ctx, championship.ID)
	require.NoError(t, err)

	// Verificações
	assert.NotEmpty(t, matches)
	// Podemos adicionar mais verificações específicas aqui
}

func TestEndToEndLeagueChampionshipCreation(t *testing.T) {
	ctx := context.Background()

	// Configurações e inicializações (mesmas do seu teste atual)
	err := config.LoadEnv()
	require.NoError(t, err)

	// Configurar o banco de dados de teste
	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	require.NotEmpty(t, dbURL, "DATABASE_URL_TEST is not set")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Limpar tabelas relevantes
	cleanupDatabase(t, pool)

	// Configurar o RabbitMQ de teste
	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL_TEST")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL_TEST is not set")

	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	require.NoError(t, err)
	defer rabbitConn.Close()

	// Inicializar repositórios e serviços (mesmo que antes)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created_test")
	require.NoError(t, err)

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo)

	rabbitMQConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created_test", matchService)
	require.NoError(t, err)

	// Iniciar consumo em uma goroutine
	go func() {
		err := rabbitMQConsumer.StartConsuming()
		require.NoError(t, err)
	}()

	userID := uuid.New()

	user := &entity.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test_league@gmail.com",
		PasswordHash: "123456",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.Create(ctx, user)

	// Criar times de teste
	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "League Team 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}
	team2 := &entity.Team{
		ID:        uuid.New(),
		Name:      "League Team 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}
	team3 := &entity.Team{
		ID:        uuid.New(),
		Name:      "League Team 3",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}
	team4 := &entity.Team{
		ID:        uuid.New(),
		Name:      "League Team 4",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}

	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team2)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team3)
	require.NoError(t, err)
	err = teamRepo.Create(ctx, team4)
	require.NoError(t, err)

	// Criar um campeonato do tipo "liga"
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test League Championship",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed, // Não é relevante para liga
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	teamIDs := []uuid.UUID{team1.ID, team2.ID, team3.ID, team4.ID}

	err = championshipService.CreateChampionship(ctx, championship, teamIDs)
	require.NoError(t, err)

	// Aguardar um tempo para o consumidor processar a mensagem
	time.Sleep(5 * time.Second)

	// Verificar se as partidas foram geradas
	matches, err := matchRepo.GetByChampionshipID(ctx, championship.ID)
	require.NoError(t, err)

	assert.NotEmpty(t, matches)

	expectedNumberOfMatches := 6
	assert.Equal(t, expectedNumberOfMatches, len(matches), "Número de partidas geradas incorreto")

	// Verificar se cada combinação única de times está presente
	combinations := make(map[string]bool)
	for _, match := range matches {
		homeID := match.HomeTeamID.String()
		awayID := match.AwayTeamID.String()
		key := homeID + "-" + awayID
		combinations[key] = true
	}

	expectedCombinations := []string{
		team1.ID.String() + "-" + team2.ID.String(),
		team1.ID.String() + "-" + team3.ID.String(),
		team1.ID.String() + "-" + team4.ID.String(),
		team2.ID.String() + "-" + team3.ID.String(),
		team2.ID.String() + "-" + team4.ID.String(),
		team3.ID.String() + "-" + team4.ID.String(),
	}

	for _, combo := range expectedCombinations {
		assert.True(t, combinations[combo], "Combinação de times não encontrada: %s", combo)
	}
}

func TestEndToEndCupChampionshipRandomDraw(t *testing.T) {
	ctx := context.Background()

	// Configurações e inicializações (mesmas do seu teste atual)
	err := config.LoadEnv()
	require.NoError(t, err)

	// Configurar o banco de dados de teste
	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	require.NotEmpty(t, dbURL, "DATABASE_URL_TEST is not set")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Limpar tabelas relevantes
	cleanupDatabase(t, pool)

	// Configurar o RabbitMQ de teste
	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL_TEST")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL_TEST is not set")

	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	require.NoError(t, err)
	defer rabbitConn.Close()

	// Inicializar repositórios e serviços (mesmo que antes)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created_test")
	require.NoError(t, err)

	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo)

	rabbitMQConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created_test", matchService)
	require.NoError(t, err)

	// Iniciar consumo em uma goroutine
	go func() {
		err := rabbitMQConsumer.StartConsuming()
		require.NoError(t, err)
	}()

	userID := uuid.New()

	user := &entity.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test_cup_random@gmail.com",
		PasswordHash: "123456",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.Create(ctx, user)

	// Criar times de teste
	teams := make([]*entity.Team, 0)
	teamIDs := make([]uuid.UUID, 0)
	numTeams := 5 // Número ímpar para testar byes

	for i := 1; i <= numTeams; i++ {
		team := &entity.Team{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("Cup Team %d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
		}
		err = teamRepo.Create(ctx, team)
		require.NoError(t, err)
		teams = append(teams, team)
		teamIDs = append(teamIDs, team.ID)
	}

	// Criar um campeonato do tipo "copa" com progressão "random_draw"
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test Cup Championship Random Draw",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionRandomDraw,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = championshipService.CreateChampionship(ctx, championship, teamIDs)
	require.NoError(t, err)

	// Aguardar um tempo para o consumidor processar a mensagem
	time.Sleep(5 * time.Second)

	// Verificar se as partidas foram geradas
	matches, err := matchRepo.GetByChampionshipID(ctx, championship.ID)
	require.NoError(t, err)

	// Verificações
	assert.NotEmpty(t, matches)

	expectedNumberOfMatches := 4
	assert.Equal(t, expectedNumberOfMatches, len(matches), "Número de partidas geradas incorreto")

	// Verificar se não há partidas duplicadas e se todas as partidas têm IDs válidos
	matchIDs := make(map[uuid.UUID]bool)
	for _, match := range matches {
		assert.False(t, matchIDs[match.ID], "Partida duplicada encontrada")
		matchIDs[match.ID] = true
	}
}

func TestEndToEndLeagueChampionshipWithEightTeams(t *testing.T) {
	ctx := context.Background()

	err := config.LoadEnv()
	require.NoError(t, err)

	// Configurar o banco de dados de teste
	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	require.NotEmpty(t, dbURL, "DATABASE_URL_TEST is not set")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Limpar tabelas relevantes
	cleanupDatabase(t, pool)

	// Configurar o RabbitMQ de teste
	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL_TEST")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL_TEST is not set")

	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	require.NoError(t, err)
	defer rabbitConn.Close()

	// Inicializar repositórios
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	// Inicializar produtor de mensagens
	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created_test")
	require.NoError(t, err)

	// Inicializar serviços
	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo)

	// Iniciar consumidor
	rabbitMQConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created_test", matchService)
	require.NoError(t, err)

	// Iniciar consumo em uma goroutine
	go func() {
		err := rabbitMQConsumer.StartConsuming()
		require.NoError(t, err)
	}()

	userID := uuid.New()

	user := &entity.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test@gmail.com",
		PasswordHash: "123456",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.Create(ctx, user)

	// Criar 8 times
	teams := make([]*entity.Team, 0)
	teamIDs := make([]uuid.UUID, 0)
	numTeams := 8

	for i := 1; i <= numTeams; i++ {
		team := &entity.Team{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("League Team %d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
		}
		err = teamRepo.Create(ctx, team)
		require.NoError(t, err)
		teams = append(teams, team)
		teamIDs = append(teamIDs, team.ID)
	}

	// Criar campeonato do tipo "liga"
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test League with Eight Teams",
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = championshipService.CreateChampionship(ctx, championship, teamIDs)
	require.NoError(t, err)

	// Aguardar processamento
	time.Sleep(3 * time.Second)

	// Verificar partidas
	matches, err := matchRepo.GetByChampionshipID(ctx, championship.ID)
	require.NoError(t, err)

	// Número esperado de partidas: n(n - 1)/2 = 28
	expectedNumberOfMatches := numTeams * (numTeams - 1) / 2
	assert.Equal(t, expectedNumberOfMatches, len(matches), "Número de partidas geradas incorreto")
}

func TestEndToEndCupChampionshipWithEightTeams(t *testing.T) {
	ctx := context.Background()

	err := config.LoadEnv()
	require.NoError(t, err)

	// Configurar o banco de dados de teste
	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	require.NotEmpty(t, dbURL, "DATABASE_URL_TEST is not set")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Limpar tabelas relevantes
	cleanupDatabase(t, pool)

	// Configurar o RabbitMQ de teste
	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL_TEST")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL_TEST is not set")

	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	require.NoError(t, err)
	defer rabbitConn.Close()

	// Inicializar repositórios
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	// Inicializar produtor de mensagens
	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created_test")
	require.NoError(t, err)

	// Inicializar serviços
	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo)

	// Iniciar consumidor
	rabbitMQConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created_test", matchService)
	require.NoError(t, err)

	// Iniciar consumo em uma goroutine
	go func() {
		err := rabbitMQConsumer.StartConsuming()
		require.NoError(t, err)
	}()

	userID := uuid.New()

	user := &entity.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test@gmail.com",
		PasswordHash: "123456",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.Create(ctx, user)

	// Criar 8 times
	teams := make([]*entity.Team, 0)
	teamIDs := make([]uuid.UUID, 0)

	numTeams := 8

	for i := 1; i <= numTeams; i++ {
		team := &entity.Team{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("Cup Team %d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
		}
		err = teamRepo.Create(ctx, team)
		require.NoError(t, err)
		teams = append(teams, team)
		teamIDs = append(teamIDs, team.ID)
	}

	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "Test Cup with Eight Teams",
		Type:             entity.ChampionshipTypeCup,
		TiebreakerMethod: entity.TiebreakerExtraTime,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = championshipService.CreateChampionship(ctx, championship, teamIDs)
	require.NoError(t, err)

	// Aguardar processamento
	time.Sleep(3 * time.Second)

	// Verificar partidas
	matches, err := matchRepo.GetByChampionshipID(ctx, championship.ID)
	require.NoError(t, err)

	// Número esperado de partidas: 7
	expectedNumberOfMatches := 7
	assert.Equal(t, expectedNumberOfMatches, len(matches), "Número de partidas geradas incorreto")
}

func TestEndToEndChampionshipCreationWithInvalidData(t *testing.T) {
	ctx := context.Background()

	err := config.LoadEnv()
	require.NoError(t, err)

	// Configurar o banco de dados de teste
	dbURL := config.GetRequiredEnv("DATABASE_URL_TEST")
	require.NotEmpty(t, dbURL, "DATABASE_URL_TEST is not set")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Limpar tabelas relevantes
	cleanupDatabase(t, pool)

	// Configurar o RabbitMQ de teste
	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL_TEST")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL_TEST is not set")

	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	require.NoError(t, err)
	defer rabbitConn.Close()

	// Inicializar repositórios
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	userRepo := repository.NewUserRepositoryPg(pool)

	// Inicializar produtor de mensagens
	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created_test")
	require.NoError(t, err)

	// Inicializar serviços
	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo)

	// Iniciar consumidor
	rabbitMQConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created_test", matchService)
	require.NoError(t, err)

	// Iniciar consumo em uma goroutine
	go func() {
		err := rabbitMQConsumer.StartConsuming()
		require.NoError(t, err)
	}()

	userID := uuid.New()

	user := &entity.User{
		ID:           userID,
		Name:         "Test User",
		Email:        "test@gmail.com",
		PasswordHash: "123456",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.Create(ctx, user)

	// Criar times válidos
	team1 := &entity.Team{
		ID:        uuid.New(),
		Name:      "Team Valid",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}
	err = teamRepo.Create(ctx, team1)
	require.NoError(t, err)

	teamIDs := []uuid.UUID{team1.ID}

	// Criar um campeonato com dados inválidos (por exemplo, nome vazio)
	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             "", // Nome inválido
		Type:             entity.ChampionshipTypeLeague,
		TiebreakerMethod: entity.TiebreakerPenalties,
		ProgressionType:  entity.ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = championshipService.CreateChampionship(ctx, championship, teamIDs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required", "Erro esperado devido a nome inválido")

	// Verificar se o campeonato não foi criado
	champ, err := championshipRepo.GetByID(ctx, championship.ID)
	require.NoError(t, err)
	assert.Nil(t, champ, "Campeonato não deveria ter sido criado")
}
