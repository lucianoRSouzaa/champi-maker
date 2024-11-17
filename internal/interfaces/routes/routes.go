package routes

import (
	"champi-maker/internal/infrastructure/config"
	"champi-maker/internal/interfaces/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	teamHandler *handler.TeamHandler,
	championshipHandler *handler.ChampionshipHandler,
	matchHandler *handler.MatchHandler,
	statisticsHandler *handler.StatisticsHandler,
) {
	router.POST("/users/register", userHandler.Register)
	router.POST("/users/login", userHandler.Login)

	jwtSecret := config.GetRequiredEnv("JWT_SECRET")
	jwtIssuer := config.GetRequiredEnv("JWT_ISSUER")
	authMiddleware := handler.AuthMiddleware(jwtSecret, jwtIssuer)

	api := router.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/users/profile", userHandler.GetProfile)
		api.PUT("/users/profile", userHandler.UpdateProfile)
		api.DELETE("/users/profile", userHandler.DeleteAccount)

		api.POST("/teams", teamHandler.CreateTeam)
		api.GET("/teams/:id", teamHandler.GetTeamByID)
		api.PUT("/teams/:id", teamHandler.UpdateTeam)
		api.DELETE("/teams/:id", teamHandler.DeleteTeam)
		api.GET("/teams", teamHandler.ListTeams)
		api.GET("/users/:user_id/teams", teamHandler.ListTeamsByUserID)

		api.POST("/championships", championshipHandler.CreateChampionship)
		api.GET("/championships/:id", championshipHandler.GetChampionshipByID)
		api.PUT("/championships/:id", championshipHandler.UpdateChampionship)
		api.DELETE("/championships/:id", championshipHandler.DeleteChampionship)
		api.GET("/championships", championshipHandler.ListChampionships)

		api.GET("/matches/:id", matchHandler.GetMatchByID)
		api.PUT("/matches/:id/result", matchHandler.UpdateMatchResult)
		api.GET("/championships/:championship_id/matches", matchHandler.ListMatchesByChampionship)

		api.POST("/championships/:id/statistics", statisticsHandler.GenerateInitialStatistics)
		api.GET("/championships/:id/statistics", statisticsHandler.GetStatisticsByChampionship)
	}
}
