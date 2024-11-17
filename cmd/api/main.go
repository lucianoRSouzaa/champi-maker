package main

import (
	"champi-maker/internal/application/service"
	security "champi-maker/internal/infrastructure/auth"
	"champi-maker/internal/infrastructure/config"
	"champi-maker/internal/infrastructure/messaging"
	"champi-maker/internal/infrastructure/repository"
	"champi-maker/internal/interfaces/handler"
	"champi-maker/internal/interfaces/routes"
	"context"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/streadway/amqp"
)

func main() {
	config.LoadEnv()

	pool := setupDatabase()
	defer pool.Close()

	userRepo := repository.NewUserRepositoryPg(pool)
	teamRepo := repository.NewTeamRepositoryPg(pool)
	championshipRepo := repository.NewChampionshipRepositoryPg(pool)
	matchRepo := repository.NewMatchRepositoryPg(pool)
	statisticsRepo := repository.NewStatisticsRepositoryPg(pool)

	jwtSecret := config.GetRequiredEnv("JWT_SECRET")
	jwtIssuer := config.GetRequiredEnv("JWT_ISSUER")
	jwtExpiry := time.Hour * 24

	tokenProvider := security.NewJWTService(jwtSecret, jwtIssuer, jwtExpiry)

	rabbitURL := config.GetRequiredEnv("RABBITMQ_URL")
	rabbitConn, err := messaging.NewRabbitMQConnection(rabbitURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao RabbitMQ: %v", err)
	}

	messagePublisher, err := messaging.NewRabbitMQPublisher(rabbitConn, "championship_created")
	if err != nil {
		log.Fatalf("Falha ao criar o MessagePublisher: %v", err)
	}

	userService := service.NewUserService(userRepo, tokenProvider)
	teamService := service.NewTeamService(teamRepo, userRepo)
	statisticsService := service.NewStatisticsService(statisticsRepo, championshipRepo, teamRepo)
	matchService := service.NewMatchService(matchRepo, championshipRepo, teamRepo, statisticsService)
	championshipService := service.NewChampionshipService(championshipRepo, teamRepo, messagePublisher)

	userHandler := handler.NewUserHandler(userService)
	teamHandler := handler.NewTeamHandler(teamService)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)
	matchHandler := handler.NewMatchHandler(matchService)
	championshipHandler := handler.NewChampionshipHandler(championshipService)

	router := gin.Default()

	routes.RegisterRoutes(router, userHandler, teamHandler, championshipHandler, matchHandler, statisticsHandler)

	go startMessageConsumer(matchService, rabbitConn)

	port := config.GetEnv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s", port)
	router.Run(":" + port)
}

func setupDatabase() *pgxpool.Pool {
	dbURL := config.GetRequiredEnv("DATABASE_URL")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	return pool
}

func startMessageConsumer(matchService service.MatchService, rabbitConn *amqp.Connection) {
	messageConsumer, err := messaging.NewRabbitMQConsumer(rabbitConn, "championship_created", matchService)
	if err != nil {
		log.Fatalf("Falha ao criar o MessageConsumer: %v", err)
	}

	err = messageConsumer.StartConsuming()
	if err != nil {
		log.Fatalf("Falha ao iniciar o consumo de mensagens: %v", err)
	}
}
