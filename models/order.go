package models

import (
	"database/sql/driver"
	"gorm.io/gorm"
)

type orderStatus string

const (
	PENDING    orderStatus = "PENDING"
	INPROGRESS orderStatus = "INPROGRESS"
	CANCELLED  orderStatus = "CANCELLED"
	COMPLETED  orderStatus = "COMPLETED"
)

func (o *orderStatus) Scan(value interface{}) error {
	*o = orderStatus(value.([]byte))
	return nil
}

func (o orderStatus) Value() (driver.Value, error) {
	return string(o), nil
}

type Order struct {
	gorm.Model
	UserProfileID uint        `gorm:"not null"`
	Status        orderStatus `sql:"type:order_status" gorm:"not null"`
	Address       string      `gorm:"not null"`
	Payment       Payment
	Products      []*Product `gorm:"many2many:product_orders"`
}

func (Order) TableName() string {
	return "orders"
}
