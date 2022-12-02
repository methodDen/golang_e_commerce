package repository

import (
	"E-Commerce/models"
	"E-Commerce/serializers/request"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

type ProductRepository interface {
	GetByStoreID(int) ([]models.Product, error)
	AddProduct(request.CreateProductInput, int) (models.Product, error)
	GetProductCategories([]int) []models.Category
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

func (p productRepository) AddProduct(newProduct request.CreateProductInput, storeID int) (models.Product, error) {

	product := models.Product{}
	err := p.DB.Where("name = ?", newProduct.Name).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		product.Name = newProduct.Name
		product.Description = newProduct.Description
		p.DB.Create(&product)
		sql := "INSERT INTO product_categories (product_id, category_id) VALUES\n"
		for i, category := range newProduct.Categories {
			if i == len(newProduct.Categories)-1 {
				sql += fmt.Sprintf("(%d, %d)\n", product.ID, category)
			} else {
				sql += fmt.Sprintf("(%d, %d),\n", product.ID, category)
			}
		}
		fmt.Println(sql)
		err := p.DB.Exec(sql).Error
		if err != nil {
			p.DB.Exec("DELETE FROM products WHERE id = ?", product.ID)
			return product, err
		}
	}
	var id int
	err = p.DB.Raw("INSERT INTO store_products (store_id, product_id, amount, price, status) VALUES(?, ?, ?, ?, ?) RETURNING store_id", storeID, product.ID, newProduct.Amount, newProduct.Price, newProduct.Status).Scan(&id).Error
	fmt.Println(err)
	return product, err
}

func (p productRepository) GetProductCategories(categoryID []int) []models.Category {
	filteredCategories := make([]models.Category, 0)
	if len(categoryID) == 0 {
		p.DB.Find(&filteredCategories)
		return filteredCategories
	}
	p.DB.Find(&filteredCategories).Where("id in ?", categoryID)
	return filteredCategories
}
