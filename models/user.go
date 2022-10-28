package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"not null; unique"`
	Role        string `gorm:"not null"`
	Password    string `gorm:"not null"`
	StoreID     *int
	UserProfile *UserProfile
}

func (User) TableName() string {
	return "users"
}

type UserProfile struct {
	gorm.Model
	FirstName   string
	LastName    string
	Country     string `gorm:"not null"`
	City        string `gorm:"not null"`
	Address     string `gorm:"not null"`
	CreditCards []CreditCard
	Orders      []Order
	UserID      uint `gorm:"not null"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
