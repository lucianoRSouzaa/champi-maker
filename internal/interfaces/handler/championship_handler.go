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

type ChampionshipHandler struct {
	service service.ChampionshipService
}

func NewChampionshipHandler(service service.ChampionshipService) *ChampionshipHandler {
	return &ChampionshipHandler{service: service}
}

type CreateChampionshipRequest struct {
	Name             string                  `json:"name" binding:"required"`
	Type             entity.ChampionshipType `json:"type" binding:"required,oneof=league cup"`
	TiebreakerMethod entity.TiebreakerMethod `json:"tiebreaker_method" validate:"required,oneof=penalties extra_time"`
	ProgressionType  entity.ProgressionType  `json:"progression_type" validate:"required,oneof=fixed random_draw"`
	TeamIDs          []uuid.UUID             `json:"team_ids" binding:"required,min=2"`
}

type UpdateChampionshipRequest struct {
	Name             string                  `json:"name" binding:"required"`
	Type             entity.ChampionshipType `json:"type" binding:"required,oneof=league cup"`
	TiebreakerMethod entity.TiebreakerMethod `json:"tiebreaker_method" validate:"required,oneof=penalties extra_time"`
	ProgressionType  entity.ProgressionType  `json:"progression_type" validate:"required,oneof=fixed random_draw"`
}

func (h *ChampionshipHandler) CreateChampionship(c *gin.Context) {
	var req CreateChampionshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	championship := &entity.Championship{
		ID:               uuid.New(),
		Name:             req.Name,
		Type:             entity.ChampionshipType(req.Type),
		TiebreakerMethod: entity.TiebreakerMethod(req.TiebreakerMethod),
		ProgressionType:  entity.ProgressionType(req.ProgressionType),
		UpdatedAt:        time.Now(),
		CreatedAt:        time.Now(),
	}

	if err := h.service.CreateChampionship(c.Request.Context(), championship, req.TeamIDs); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusCreated, championship)
}

func (h *ChampionshipHandler) GetChampionshipByID(c *gin.Context) {
	idParam := c.Param("id")
	championshipID, err := uuid.Parse(idParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "invalid championship ID")
		return
	}

	championship, err := h.service.GetChampionshipByID(c.Request.Context(), championshipID)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if championship == nil {
		web.RespondWithError(c, http.StatusNotFound, "championship not found")
		return
	}

	web.RespondWithJSON(c, http.StatusOK, championship)
}

func (h *ChampionshipHandler) UpdateChampionship(c *gin.Context) {
	idParam := c.Param("id")
	championshipID, err := uuid.Parse(idParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "invalid championship ID")
		return
	}

	var req UpdateChampionshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	championship := &entity.Championship{
		ID:               championshipID,
		Name:             req.Name,
		Type:             entity.ChampionshipType(req.Type),
		TiebreakerMethod: entity.TiebreakerMethod(req.TiebreakerMethod),
		ProgressionType:  entity.ProgressionType(req.ProgressionType),
	}

	if err := h.service.UpdateChampionship(c.Request.Context(), championship); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, championship)
}

func (h *ChampionshipHandler) DeleteChampionship(c *gin.Context) {
	idParam := c.Param("id")
	championshipID, err := uuid.Parse(idParam)
	if err != nil {
		web.RespondWithError(c, http.StatusBadRequest, "invalid championship ID")
		return
	}

	if err := h.service.DeleteChampionship(c.Request.Context(), championshipID); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, gin.H{"message": "championship deleted successfully"})
}

func (h *ChampionshipHandler) ListChampionships(c *gin.Context) {
	championships, err := h.service.ListChampionships(c.Request.Context())
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, championships)
}
