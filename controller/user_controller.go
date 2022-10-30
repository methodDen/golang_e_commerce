package controller

import (
	"E-Commerce/common"
	"E-Commerce/models"
	"E-Commerce/repository"
	"E-Commerce/serializers/request"
	"fmt"
	"github.com/casbin/casbin/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController : represent the user's controller contract
type UserController interface {
	AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc
	GetUser(*gin.Context)
	SignInUser(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}

type userController struct {
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepository
}

//NewUserController -> returns new user controller
func NewUserController(userRepo repository.UserRepository, storeRepo repository.StoreRepository) UserController {
	return userController{
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
}

func (h userController) GetUser(ctx *gin.Context) {
	id := ctx.Param("user")
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userRepo.GetUser(intID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}

func (h userController) SignInUser(ctx *gin.Context) {
	var input request.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	dbUser, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "No Such User Found"})
		return
	}
	if isTrue := common.ComparePassword(dbUser.Password, input.Password); isTrue {
		fmt.Println("user before", dbUser.ID)
		token := common.GenerateToken(dbUser.ID)
		ctx.JSON(http.StatusOK, gin.H{"msg": "Successfully SignIN", "token": token})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Password not matched"})
	return

}

func (h userController) AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input request.CreateUserInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userExists, err := h.userRepo.UserExists(input.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if userExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
		common.HashPassword(&input.Password)
		var user models.User
		if input.Role == "STOREADMIN" {
			storeId := input.StoreID
			fmt.Println(storeId)
			storeExists, err := h.storeRepo.StoreExists(storeId)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !storeExists {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Store does not exists"})
				return
			}
			var p *int
			p = &storeId
			user = models.User{Email: input.Email, Role: input.Role, Password: input.Password, StoreID: p}
		} else {
			user = models.User{Email: input.Email, Role: input.Role, Password: input.Password, StoreID: nil}
		}
		db := models.GetDB()
		db.Create(&user)
		db.Create(&models.UserProfile{UserID: user.ID})
		enforcer.AddGroupingPolicy(fmt.Sprint(user.ID), user.Role)
		user.Password = ""
		ctx.JSON(http.StatusOK, user)

	}
}

func (h userController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("user")
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user.ID = uint(intID)
	user, err = h.userRepo.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}

func (h userController) DeleteUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("user")
	intID, _ := strconv.Atoi(id)
	user.ID = uint(intID)
	user, err := h.userRepo.DeleteUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, user)

}
