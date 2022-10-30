package controller

import (
	"E-Commerce/common"
	"E-Commerce/models"
	"E-Commerce/repository"
	"E-Commerce/serializers/request"
	"E-Commerce/serializers/response"
	"fmt"
	"github.com/casbin/casbin/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController : represent the user's controller contract
type UserController interface {
	AddUser(enforcer *casbin.Enforcer) gin.HandlerFunc
	GetUserProfile(*gin.Context)
	SignInUser(*gin.Context)
	UpdateUserProfile(*gin.Context)
	//UpdateUserCredentials(*gin.Context)
}

type userController struct {
	userRepo        repository.UserRepository
	storeRepo       repository.StoreRepository
	userProfileRepo repository.UserProfileRepository
}

func NewUserController(
	userRepo repository.UserRepository,
	storeRepo repository.StoreRepository,
	userProfileRepo repository.UserProfileRepository,
) UserController {
	return userController{
		userRepo:        userRepo,
		storeRepo:       storeRepo,
		userProfileRepo: userProfileRepo,
	}
}

func (h userController) GetUserProfile(ctx *gin.Context) {
	id, _ := ctx.Get("userID")
	intID := int(id.(float64))
	userProfile, err := h.userProfileRepo.GetUserProfile(intID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := response.UserProfileResponse{
		ID:        userProfile.ID,
		UserID:    userProfile.UserID,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Country:   userProfile.Country,
		City:      userProfile.City,
		Address:   userProfile.Address,
	}
	ctx.JSON(http.StatusOK, resp)

}

func (h userController) UpdateUserProfile(context *gin.Context) {
	var input request.UpdateUserProfileInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := context.Get("userID")
	intID := int(id.(float64))
	userProfile, err := h.userProfileRepo.GetUserProfile(intID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userProfile, err = h.userProfileRepo.UpdateUserProfile(userProfile, input)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := response.UserProfileResponse{
		ID:        userProfile.ID,
		UserID:    userProfile.UserID,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Country:   userProfile.Country,
		City:      userProfile.City,
		Address:   userProfile.Address,
	}
	context.JSON(http.StatusOK, resp)
}

func (h userController) SignInUser(ctx *gin.Context) {
	var input request.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	dbUser, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "No Such User Found"})
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
		userExists, err := h.userRepo.UserExistsByEmail(input.Email)
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
