package handlers

import (
	"e-wallet-api/internal/config"
	"e-wallet-api/internal/services"
	"e-wallet-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid")
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	response := map[string]interface{}{
		"user_id":      user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"phone_number": user.PhoneNumber,
		"address":      user.Address,
		"created_date": user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.SuccessResponse(c, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid")
		return
	}

	user, err := h.authService.Login(req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, h.config.JWT.Secret)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Gagal membuat token")
		return
	}

	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	utils.SuccessResponse(c, response)
}
