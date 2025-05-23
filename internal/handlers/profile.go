package handlers

import (
	"e-wallet-api/internal/services"
	"e-wallet-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	walletService *services.WalletService
}

func NewProfileHandler(walletService *services.WalletService) *ProfileHandler {
	return &ProfileHandler{
		walletService: walletService,
	}
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User tidak terautentikasi")
		return
	}

	var req services.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid")
		return
	}

	user, err := h.walletService.UpdateProfile(userID.(uuid.UUID), req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	response := map[string]interface{}{
		"user_id":      user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"address":      user.Address,
		"updated_date": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.SuccessResponse(c, response)
}
