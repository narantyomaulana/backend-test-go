package services

import (
	"errors"

	"e-wallet-api/internal/database"
	"e-wallet-api/internal/models"
	"e-wallet-api/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {

	if !utils.ValidateName(req.FirstName) || !utils.ValidateName(req.LastName) {
		return nil, errors.New("nama firstname atau lastname tidak boleh kosong")
	}

	if !utils.ValidatePhoneNumber(req.PhoneNumber) {
		return nil, errors.New("nomor telepon tidak valid")
	}

	if !utils.ValidatePIN(req.PIN) {
		return nil, errors.New("PIN harus 6 digit angka")
	}

	var existingUser models.User
	if err := database.GetDB().Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		return nil, errors.New("nomor telepon sudah terdaftar")
	}

	// Hash PIN
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi PIN")
	}

	// Create user
	user := models.User{
		ID:          uuid.New(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		PIN:         string(hashedPIN),
		Balance:     0,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		return nil, errors.New("gagal membuat akun")
	}

	return &user, nil
}

func (s *AuthService) Login(req LoginRequest) (*models.User, error) {
	var user models.User
	if err := database.GetDB().Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("phone number dan PIN tidak cocok")
		}
		return nil, errors.New("gagal login")
	}

	// Verify PIN
	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(req.PIN)); err != nil {
		return nil, errors.New("phone number dan PIN tidak cocok")
	}

	return &user, nil
}

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PIN         string `json:"pin" binding:"required"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PIN         string `json:"pin" binding:"required"`
}
