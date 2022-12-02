package main

import (
	"E-Commerce/models"
	"E-Commerce/route"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	err := db.Migrator().DropTable(
		&models.Store{},
		&models.UserProfile{},
		&models.User{},
		&models.StoreProduct{},
		&models.CreditCard{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.ProductOrder{},
		&models.Payment{})
	if err != nil {
		log.Println(err)
	}
	if err := db.AutoMigrate(
		&models.Store{},
		&models.User{},
		&models.UserProfile{},
		&models.Product{},
		&models.CreditCard{},
		&models.Category{},
		&models.Order{},
		&models.StoreProduct{},
		&models.ProductOrder{},
		&models.Payment{},
	); err != nil {
		log.Println(err)
	}
}

func main() {
	db := models.Init()
	//Migrate(db)
	route.SetupRoutes(db)
}
