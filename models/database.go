package models

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	// read configs from .env
	db_user := viper.Get("DB_USER").(string)
	db_password := viper.Get("DB_PASSWORD").(string)
	db_name := viper.Get("DB_NAME").(string)
	db_host := viper.Get("DB_HOST").(string)
	db_port := viper.Get("DB_PORT").(string)
	url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Almaty",
		db_host,
		db_user,
		db_password,
		db_name,
		db_port,
	)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database!")
	}
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
