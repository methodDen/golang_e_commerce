package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserProfileID uint   `gorm:"not null"`
	Status        string `gorm:"not null"`
	Address       string `gorm:"not null"`
	Payment       Payment
	Products      []*Product `gorm:"many2many:product_orders"`
	StoreID       int
	Store         Store
}

func (Order) TableName() string {
	return "orders"
}
