package controller

import (
	"E-Commerce/repository"
	"E-Commerce/serializers/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StoreController interface {
	GetProductsByStoreID(*gin.Context)
}

type storeController struct {
	storeRepo   repository.StoreRepository
	productRepo repository.ProductRepository
}

func (s storeController) GetProductsByStoreID(context *gin.Context) {
	studioID := context.Param("id")
	intID, err := strconv.Atoi(studioID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	products, err := s.productRepo.GetByStoreID(intID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var resp []response.ProductResponse
	for _, product := range products {
		serializer := response.ProductResponse{ID: product.ID, Name: product.Name, Description: product.Description}
		resp = append(resp, serializer)
	}
	context.JSON(http.StatusOK, resp)
}

//NewUserController -> returns new user controller
func NewStoreController(storeRepo repository.StoreRepository, productRepo repository.ProductRepository) StoreController {
	return storeController{
		storeRepo:   storeRepo,
		productRepo: productRepo,
	}
}
