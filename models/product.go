package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Categories  []*Category `gorm:"many2many:product_categories"`
	Stores      []*Store    `gorm:"many2many:store_products"`
	Orders      []*Order    `gorm:"many2many:product_orders"`
}

func (Product) TableName() string {
	return "products"
}

type ProductOrder struct {
	ProductID uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"primaryKey"`
	Quantity  uint `gorm:"not null"`
	Price     uint `gorm:"not null"`
}

func (ProductOrder) TableName() string {
	return "product_orders"
}
