package repository

import (
	"E-Commerce/models"
	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

type ProductRepository interface {
	GetByStoreID(int) ([]models.Product, error)
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return productRepository{
		DB: db,
	}
}

func (p productRepository) GetByStoreID(storeID int) (products []models.Product, err error) {
	store := models.Store{}
	p.DB.First(&store, "id = ?", storeID)
	return products, p.DB.Model(&store).Association("Products").Find(&products)
}
