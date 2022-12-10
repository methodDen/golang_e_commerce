package repository

import (
	"E-Commerce/models"
	"errors"
	"gorm.io/gorm"
)

type ordersRepo struct {
	DB *gorm.DB
}

type OrdersRepo interface {
	GetStoreOrders(int) ([]models.Order, error)
	DoneStoreOrder(int) (models.Order, error)
	TakeStoreOrder(int) (models.Order, error)
	RejectStoreOrder(int) (models.Order, error)
	AddStoreOrder(int, uint, string, []int) (models.Order, error)
	GetStoreOrder(int, int) (models.Order, error)
}

func NewOrdersRepo(db *gorm.DB) OrdersRepo {
	return ordersRepo{
		DB: db,
	}
}

func (r ordersRepo) GetStoreOrders(storeID int) ([]models.Order, error) {
	var orders = make([]models.Order, 0)
	err := r.DB.Where("storeID = ?", storeID).Find(&orders)
	if err != nil {
		return orders, err.Error
	}
	return orders, nil
}

func (r ordersRepo) GetStoreOrder(storeID int, orderID int) (models.Order, error) {
	var order models.Order
	err := r.DB.First(&order, orderID)
	if err != nil {
		return models.Order{}, err.Error
	}
	return order, nil
}

func (r ordersRepo) DoneStoreOrder(productID int) (models.Order, error) {
	var order models.Order
	err := r.DB.First(&order, productID)
	if err != nil {
		return order, err.Error
	}
	order.Status = "DONE"
	return order, nil
}

func (r ordersRepo) RejectStoreOrder(productID int) (models.Order, error) {
	var order models.Order
	err := r.DB.First(&order, productID)
	if err != nil {
		return order, err.Error
	}
	order.Status = "REJECTED"
	return order, nil
}

func (r ordersRepo) TakeStoreOrder(productID int) (models.Order, error) {
	var order models.Order
	if err := r.DB.First(&order, productID); err != nil {
		return order, err.Error
	}
	order.Status = "DELIVERY"
	return order, nil
}

func (r ordersRepo) AddStoreOrder(storeID int, userID uint, address string, products []int) (models.Order, error) {
	productsInDB := make([]models.Product, 0)
	err := r.DB.Where("name IN ?", products).Find(&productsInDB)
	if err != nil {
		return models.Order{}, err.Error
	}
	if len(products) != len(productsInDB) {
		return models.Order{}, errors.New("invalid product ids")
	}
	order := models.Order{}
	order.StoreID = storeID
	order.Status = "PENDING"
	order.Address = address
	order.UserProfileID = userID
	err = r.DB.Create(&order)
	return order, err.Error
}
