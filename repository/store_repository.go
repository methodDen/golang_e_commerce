package repository

import (
	"E-Commerce/models"
	"gorm.io/gorm"
)

type storeRepository struct {
	DB *gorm.DB
}

type StoreRepository interface {
	StoreExists(int) (bool, error)
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return storeRepository{
		DB: db,
	}
}

func (s storeRepository) StoreExists(id int) (bool, error) {
	var store models.Store
	err := s.DB.First(&store, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
