package handlers

import (
	"e-wallet-api/internal/services"
	"e-wallet-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	walletService *services.WalletService
}

func NewPaymentHandler(walletService *services.WalletService) *PaymentHandler {
	return &PaymentHandler{
		walletService: walletService,
	}
}

func (h *PaymentHandler) Payment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User tidak terautentikasi")
		return
	}

	var req services.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid")
		return
	}

	payment, err := h.walletService.Payment(userID.(uuid.UUID), req.Amount, req.Remarks)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	response := map[string]interface{}{
		"payment_id":     payment.ID,
		"amount":         payment.Amount,
		"remarks":        payment.Remarks,
		"balance_before": payment.BalanceBefore,
		"balance_after":  payment.BalanceAfter,
		"created_date":   payment.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.SuccessResponse(c, response)
}
