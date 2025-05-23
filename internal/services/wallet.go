package services

import (
	"errors"
	"log"

	"e-wallet-api/internal/database"
	"e-wallet-api/internal/models"
	"e-wallet-api/internal/utils"
	"e-wallet-api/pkg/rabbitmq"

	"github.com/google/uuid"
)

type WalletService struct {
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewWalletService(rabbitMQ *rabbitmq.RabbitMQ) *WalletService {
	return &WalletService{
		rabbitMQ: rabbitMQ,
	}
}

func (s *WalletService) TopUp(userID uuid.UUID, amount int64) (*models.TopUp, error) {
	if !utils.ValidateAmount(amount) {
		return nil, errors.New("jumlah top up tidak valid")
	}

	var user models.User
	if err := database.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Begin transaction
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	balanceBefore := user.Balance
	balanceAfter := balanceBefore + amount

	// Update user balance
	if err := tx.Model(&user).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal mengupdate saldo")
	}

	// Create top up record
	topUp := models.TopUp{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
	}

	if err := tx.Create(&topUp).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat record top up")
	}

	// Create transaction record
	transaction := models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		TransactionType: "CREDIT",
		Amount:          amount,
		Remarks:         "Top Up",
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Status:          "SUCCESS",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat record transaksi")
	}

	tx.Commit()
	return &topUp, nil
}

func (s *WalletService) Payment(userID uuid.UUID, amount int64, remarks string) (*models.Payment, error) {
	if !utils.ValidateAmount(amount) {
		return nil, errors.New("jumlah pembayaran tidak valid")
	}

	if remarks == "" {
		return nil, errors.New("keterangan pembayaran tidak boleh kosong")
	}

	var user models.User
	if err := database.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.Balance < amount {
		return nil, errors.New("saldo tidak cukup")
	}

	// Begin transaction
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	balanceBefore := user.Balance
	balanceAfter := balanceBefore - amount

	// Update user balance
	if err := tx.Model(&user).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal mengupdate saldo")
	}

	// Create payment record
	payment := models.Payment{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat record pembayaran")
	}

	// Create transaction record
	transaction := models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		TransactionType: "DEBIT",
		Amount:          amount,
		Remarks:         remarks,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Status:          "SUCCESS",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat record transaksi")
	}

	tx.Commit()
	return &payment, nil
}

func (s *WalletService) InitiateTransfer(fromUserID, toUserID uuid.UUID, amount int64, remarks string) (*models.Transfer, error) {
	if !utils.ValidateAmount(amount) {
		return nil, errors.New("jumlah transfer tidak valid")
	}

	if remarks == "" {
		return nil, errors.New("keterangan transfer tidak boleh kosong")
	}

	if fromUserID == toUserID {
		return nil, errors.New("tidak dapat transfer ke diri sendiri")
	}

	// Check if both users exist
	var fromUser, toUser models.User
	if err := database.GetDB().Where("id = ?", fromUserID).First(&fromUser).Error; err != nil {
		return nil, errors.New("user pengirim tidak ditemukan")
	}

	if err := database.GetDB().Where("id = ?", toUserID).First(&toUser).Error; err != nil {
		return nil, errors.New("user penerima tidak ditemukan")
	}

	if fromUser.Balance < amount {
		return nil, errors.New("saldo tidak cukup")
	}

	// Begin transaction for sender
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	balanceBefore := fromUser.Balance
	balanceAfter := balanceBefore - amount

	// Update sender balance
	if err := tx.Model(&fromUser).Update("balance", balanceAfter).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal mengupdate saldo pengirim")
	}

	// Create transfer record
	transfer := models.Transfer{
		ID:            uuid.New(),
		FromUserID:    fromUserID,
		ToUserID:      toUserID,
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        "PENDING",
	}

	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat record transfer")
	}

	// Create transaction record for sender
	transaction := models.Transaction{
		ID:              uuid.New(),
		UserID:          fromUserID,
		TransactionType: "DEBIT",
		Amount:          amount,
		Remarks:         "Transfer ke " + toUser.FirstName + " " + toUser.LastName + " - " + remarks,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Status:          "SUCCESS",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat record transaksi")
	}

	tx.Commit()

	// Send to queue for background processing
	message := rabbitmq.TransferMessage{
		TransferID: transfer.ID.String(),
		FromUserID: fromUserID.String(),
		ToUserID:   toUserID.String(),
		Amount:     amount,
		Remarks:    remarks,
	}

	if err := s.rabbitMQ.PublishMessage("transfer_queue", message); err != nil {
		log.Printf("Failed to publish transfer message: %v", err)
		// Note: Transfer already deducted from sender, will be processed by background worker
	}

	return &transfer, nil
}

func (s *WalletService) ProcessTransfer(message rabbitmq.TransferMessage) error {
	transferID, err := uuid.Parse(message.TransferID)
	if err != nil {
		return err
	}

	toUserID, err := uuid.Parse(message.ToUserID)
	if err != nil {
		return err
	}

	fromUserID, err := uuid.Parse(message.FromUserID)
	if err != nil {
		return err
	}

	// Get transfer record
	var transfer models.Transfer
	if err := database.GetDB().Where("id = ?", transferID).First(&transfer).Error; err != nil {
		return err
	}

	// Get recipient user
	var toUser models.User
	if err := database.GetDB().Where("id = ?", toUserID).First(&toUser).Error; err != nil {
		return err
	}

	var fromUser models.User
	if err := database.GetDB().Where("id = ?", fromUserID).First(&fromUser).Error; err != nil {
		return err
	}

	// Begin transaction
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update recipient balance
	newBalance := toUser.Balance + message.Amount
	if err := tx.Model(&toUser).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update transfer status
	if err := tx.Model(&transfer).Update("status", "SUCCESS").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create transaction record for recipient
	transaction := models.Transaction{
		ID:              uuid.New(),
		UserID:          toUserID,
		TransactionType: "CREDIT",
		Amount:          message.Amount,
		Remarks:         "Transfer dari " + fromUser.FirstName + " " + fromUser.LastName + " - " + message.Remarks,
		BalanceBefore:   toUser.Balance,
		BalanceAfter:    newBalance,
		Status:          "SUCCESS",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	log.Printf("Transfer processed successfully: %s", transferID)
	return nil
}

func (s *WalletService) GetTransactions(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, errors.New("gagal mengambil data transaksi")
	}
	return transactions, nil
}

func (s *WalletService) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*models.User, error) {
	if !utils.ValidateName(req.FirstName) || !utils.ValidateName(req.LastName) {
		return nil, errors.New("nama firstname atau lastname tidak boleh kosong")
	}

	var user models.User
	if err := database.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Update user profile
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Address = req.Address

	if err := database.GetDB().Save(&user).Error; err != nil {
		return nil, errors.New("gagal mengupdate profil")
	}

	return &user, nil
}

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

type PaymentRequest struct {
	Amount  int64  `json:"amount" binding:"required"`
	Remarks string `json:"remarks" binding:"required"`
}

type TransferRequest struct {
	TargetUser string `json:"target_user" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
	Remarks    string `json:"remarks" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
}
