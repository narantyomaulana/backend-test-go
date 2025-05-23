package handlers

import (
	"e-wallet-api/internal/services"
	"e-wallet-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TopUpHandler struct {
	walletService *services.WalletService
}

func NewTopUpHandler(walletService *services.WalletService) *TopUpHandler {
	return &TopUpHandler{
		walletService: walletService,
	}
}

func (h *TopUpHandler) TopUp(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User tidak terautentikasi")
		return
	}

	var req services.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid")
		return
	}

	topUp, err := h.walletService.TopUp(userID.(uuid.UUID), req.Amount)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	response := map[string]interface{}{
		"top_up_id":      topUp.ID,
		"amount_top_up":  topUp.Amount,
		"balance_before": topUp.BalanceBefore,
		"balance_after":  topUp.BalanceAfter,
		"created_date":   topUp.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.SuccessResponse(c, response)
}
