package models

import (
	"gorm.io/gorm"
	"time"
)

type CreditCard struct {
	gorm.Model
	CVCode            string    `gorm:"not null; size:3"`
	DateOfManufacture time.Time `gorm:"not null"`
	DateOfExpiry      time.Time `gorm:"not null"`
	Amount            uint      `gorm:"not null"`
	UserProfileID     uint      `gorm:"not null"`
	Payments          []Payment
}

func (CreditCard) TableName() string {
	return "credit_cards"
}
