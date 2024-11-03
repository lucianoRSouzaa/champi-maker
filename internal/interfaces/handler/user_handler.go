package handler

import (
	"champi-maker/internal/application/service"
	"champi-maker/internal/domain/entity"
	"champi-maker/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := &entity.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.userService.RegisterUser(c.Request.Context(), user, req.Password); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user.PasswordHash = ""

	web.RespondWithJSON(c, http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.userService.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		web.RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
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

	user, err := h.userService.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	user.PasswordHash = ""

	web.RespondWithJSON(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := &entity.User{
		ID: uid,
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		user.PasswordHash = *req.Password
	}

	if err := h.userService.UpdateUser(c.Request.Context(), user); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Perfil atualizado com sucesso"})
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
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

	if err := h.userService.DeleteUser(c.Request.Context(), uid); err != nil {
		web.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Conta excluída com sucesso"})
}
