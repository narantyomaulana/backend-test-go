package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"user_id"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
	PhoneNumber string    `gorm:"unique;not null" json:"phone_number"`
	Address     string    `gorm:"not null" json:"address"`
	PIN         string    `gorm:"not null" json:"-"`
	Balance     int64     `gorm:"default:0" json:"balance"`
	CreatedAt   time.Time `json:"created_date"`
	UpdatedAt   time.Time `json:"updated_date"`
}

type TopUp struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"top_up_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount        int64     `gorm:"not null" json:"amount_top_up"`
	BalanceBefore int64     `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64     `gorm:"not null" json:"balance_after"`
	CreatedAt     time.Time `json:"created_date"`
	User          User      `gorm:"foreignKey:UserID"`
}

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"payment_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount        int64     `gorm:"not null" json:"amount"`
	Remarks       string    `gorm:"not null" json:"remarks"`
	BalanceBefore int64     `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64     `gorm:"not null" json:"balance_after"`
	CreatedAt     time.Time `json:"created_date"`
	User          User      `gorm:"foreignKey:UserID"`
}

type Transfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"transfer_id"`
	FromUserID    uuid.UUID `gorm:"type:uuid;not null" json:"from_user_id"`
	ToUserID      uuid.UUID `gorm:"type:uuid;not null" json:"to_user_id"`
	Amount        int64     `gorm:"not null" json:"amount"`
	Remarks       string    `gorm:"not null" json:"remarks"`
	BalanceBefore int64     `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64     `gorm:"not null" json:"balance_after"`
	Status        string    `gorm:"default:'PENDING'" json:"status"`
	CreatedAt     time.Time `json:"created_date"`
	FromUser      User      `gorm:"foreignKey:FromUserID"`
	ToUser        User      `gorm:"foreignKey:ToUserID"`
}

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"transaction_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TransactionType string    `gorm:"not null" json:"transaction_type"` // DEBIT or CREDIT
	Amount          int64     `gorm:"not null" json:"amount"`
	Remarks         string    `gorm:"not null" json:"remarks"`
	BalanceBefore   int64     `gorm:"not null" json:"balance_before"`
	BalanceAfter    int64     `gorm:"not null" json:"balance_after"`
	Status          string    `gorm:"default:'SUCCESS'" json:"status"`
	CreatedAt       time.Time `json:"created_date"`
	User            User      `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (t *TopUp) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (t *Transfer) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
