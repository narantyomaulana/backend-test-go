package handlers

import (
	"e-wallet-api/internal/services"
	"e-wallet-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransferHandler struct {
	walletService *services.WalletService
}

func NewTransferHandler(walletService *services.WalletService) *TransferHandler {
	return &TransferHandler{
		walletService: walletService,
	}
}

func (h *TransferHandler) Transfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User tidak terautentikasi")
		return
	}

	var req services.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid")
		return
	}

	targetUserID, err := uuid.Parse(req.TargetUser)
	if err != nil {
		utils.ValidationErrorResponse(c, "Target user ID tidak valid")
		return
	}

	transfer, err := h.walletService.InitiateTransfer(userID.(uuid.UUID), targetUserID, req.Amount, req.Remarks)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	response := map[string]interface{}{
		"transfer_id":    transfer.ID,
		"amount":         transfer.Amount,
		"remarks":        transfer.Remarks,
		"balance_before": transfer.BalanceBefore,
		"balance_after":  transfer.BalanceAfter,
		"created_date":   transfer.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.SuccessResponse(c, response)
}
