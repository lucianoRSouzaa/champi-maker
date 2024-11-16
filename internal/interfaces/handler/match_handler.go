package handler

import (
	"champi-maker/internal/application/service"
	"champi-maker/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchHandler struct {
	matchService service.MatchService
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
	}
}

type UpdateMatchResultRequest struct {
	ScoreHome          int  `json:"score_home" binding:"required"`
	ScoreAway          int  `json:"score_away" binding:"required"`
	HasExtraTime       bool `json:"has_extra_time"`
	ScoreHomeExtraTime int  `json:"score_home_extra_time,omitempty"`
	ScoreAwayExtraTime int  `json:"score_away_extra_time,omitempty"`
	HasPenalties       bool `json:"has_penalties"`
	ScoreHomePenalties int  `json:"score_home_penalties,omitempty"`
	ScoreAwayPenalties int  `json:"score_away_penalties,omitempty"`
}

func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	matchIDParam := c.Param("id")
	matchID, err := uuid.Parse(matchIDParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID da partida inválido")
		return
	}

	match, err := h.matchService.GetMatchByID(c.Request.Context(), matchID)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, match)
}

func (h *MatchHandler) ListMatchesByChampionship(c *gin.Context) {
	championshipIDParam := c.Param("championship_id")
	championshipID, err := uuid.Parse(championshipIDParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID do campeonato inválido")
		return
	}

	matches, err := h.matchService.ListMatchesByChampionship(c.Request.Context(), championshipID)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, matches)
}

func (h *MatchHandler) UpdateMatchResult(c *gin.Context) {
	matchIDParam := c.Param("id")
	matchID, err := uuid.Parse(matchIDParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID da partida inválido")
		return
	}

	var req UpdateMatchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Criar a estrutura MatchResultUpdate com base na requisição
	resultUpdate := service.MatchResultUpdate{
		ScoreHome:          req.ScoreHome,
		ScoreAway:          req.ScoreAway,
		HasExtraTime:       req.HasExtraTime,
		ScoreHomeExtraTime: req.ScoreHomeExtraTime,
		ScoreAwayExtraTime: req.ScoreAwayExtraTime,
		HasPenalties:       req.HasPenalties,
		ScoreHomePenalties: req.ScoreHomePenalties,
		ScoreAwayPenalties: req.ScoreAwayPenalties,
	}

	if err := h.matchService.UpdateMatchResult(c.Request.Context(), matchID, resultUpdate); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Resultado da partida atualizado com sucesso"})
}
