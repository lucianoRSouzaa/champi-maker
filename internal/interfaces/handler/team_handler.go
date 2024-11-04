package handler

import (
	"champi-maker/internal/application/service"
	"champi-maker/internal/domain/entity"
	"champi-maker/pkg/web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		web.RespondWithError(c, http.StatusUnauthorized, "ID de usuário não encontrado no contexto")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		web.RespondWithError(c, http.StatusUnauthorized, "ID de usuário inválido")
		return
	}

	team := &entity.Team{
		ID:        uuid.New(),
		Name:      req.Name,
		UserID:    uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.teamService.CreateTeam(c.Request.Context(), team); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusCreated, team)
}

func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	idParam := c.Param("id")
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID de time inválido")
		return
	}

	team, err := h.teamService.GetTeamByID(c.Request.Context(), teamID)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if team == nil {
		web.RespondWithError(c, http.StatusNotFound, "time não encontrado")
		return
	}

	web.RespondWithJSON(c, http.StatusOK, team)
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	idParam := c.Param("id")
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID de time inválido")
		return
	}

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	team := &entity.Team{
		ID:        teamID,
		Name:      req.Name,
		UpdatedAt: time.Now(),
	}

	if err := h.teamService.UpdateTeam(c.Request.Context(), team); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, team)
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	idParam := c.Param("id")
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID de time inválido")
		return
	}

	if err := h.teamService.DeleteTeam(c.Request.Context(), teamID); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, gin.H{"message": "time excluído com sucesso"})
}

func (h *TeamHandler) ListTeams(c *gin.Context) {
	teams, err := h.teamService.ListTeams(c.Request.Context())
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, teams)
}

func (h *TeamHandler) ListTeamsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "ID de usuário inválido")
		return
	}

	teams, err := h.teamService.ListTeamsByUserID(c.Request.Context(), userID)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, teams)
}
