package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string     `gorm:"not null"`
	Products []*Product `gorm:"many2many:product_categories"`
}

func (Category) TableName() string {
	return "categories"
}
