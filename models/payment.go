package models

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	OrderID         uint      `gorm:"not null"`
	TotalPrice      uint      `gorm:"not null"`
	CreditCardID    uint      `gorm:"not null"`
	TransactionDate time.Time `gorm:"not null"`
}

func (Payment) TableName() string {
	return "payments"
}
