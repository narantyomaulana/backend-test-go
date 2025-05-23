package handlers

import (
	"e-wallet-api/internal/services"
	"e-wallet-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	walletService *services.WalletService
}

func NewTransactionHandler(walletService *services.WalletService) *TransactionHandler {
	return &TransactionHandler{
		walletService: walletService,
	}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User tidak terautentikasi")
		return
	}

	transactions, err := h.walletService.GetTransactions(userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	var response []map[string]interface{}
	for _, transaction := range transactions {
		response = append(response, map[string]interface{}{
			"transaction_id":   transaction.ID,
			"status":           transaction.Status,
			"user_id":          transaction.UserID,
			"transaction_type": transaction.TransactionType,
			"amount":           transaction.Amount,
			"remarks":          transaction.Remarks,
			"balance_before":   transaction.BalanceBefore,
			"balance_after":    transaction.BalanceAfter,
			"created_date":     transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.SuccessResponse(c, response)
}
