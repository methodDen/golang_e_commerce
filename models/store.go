package models

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Country     string `gorm:"not null"`
	City        string `gorm:"not null"`
	Address     string `gorm:"not null"`
	Description string
	PhoneNumber string
	Website     string
	Users       []User
	Products    []*Product `gorm:"many2many:store_products"`
}

func (Store) TableName() string {
	return "stores"
}

type StoreProduct struct {
	StoreID   uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Amount    uint `gorm:"not null"`
}

func (StoreProduct) TableName() string {
	return "store_products"
}
