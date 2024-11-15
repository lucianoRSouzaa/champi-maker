package handler

import (
	"champi-maker/internal/application/service"
	"champi-maker/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StatisticsHandler struct {
	statisticsService service.StatisticsService
}

func NewStatisticsHandler(statisticsService service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

type GenerateStatisticsRequest struct {
	TeamIDs []uuid.UUID `json:"team_ids" binding:"required,min=1"`
}

func (h *StatisticsHandler) GenerateInitialStatistics(c *gin.Context) {
	championshipIDParam := c.Param("id")
	championshipID, err := uuid.Parse(championshipIDParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID de campeonato inválido")
		return
	}

	var req GenerateStatisticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.statisticsService.GenerateInitialStatistics(c.Request.Context(), championshipID, req.TeamIDs)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusCreated, gin.H{"message": "Estatísticas iniciais geradas com sucesso"})
}

func (h *StatisticsHandler) GetStatisticsByChampionship(c *gin.Context) {
	championshipIDParam := c.Param("id")
	championshipID, err := uuid.Parse(championshipIDParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID de campeonato inválido")
		return
	}

	statsList, err := h.statisticsService.GetStatisticsByChampionship(c.Request.Context(), championshipID)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, statsList)
}
