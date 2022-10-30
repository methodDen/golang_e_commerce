package route

import (
	"E-Commerce/controller"
	"E-Commerce/repository"
	//"E-Commerce/controller"
	"E-Commerce/middleware"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//SetupRoutes : all the routes are defined here

func SetupEnforcerPolicies(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//add policy for client and story_admin
	if hasPolicy := enforcer.HasPolicy("CLIENT", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("CLIENT", "report", "read")
	}
	if hasPolicy := enforcer.HasPolicy("STOREADMIN", "report", "write"); !hasPolicy {
		enforcer.AddPolicy("STOREADMIN", "report", "write")
	}
	if hasPolicy := enforcer.HasPolicy("CLIENT", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("STOREADMIN", "report", "read")
	}

	// add policy for supermoderator

	return enforcer
}

func SetupRoutes(db *gorm.DB) {
	httpRouter := gin.Default()

	// Initialize  repositories

	userRepository := repository.NewUserRepository(db)
	storeRepository := repository.NewStoreRepository(db)
	productRepository := repository.NewProductRepository(db)
	userProfileRepository := repository.NewUserProfileRepository(db)
	// Initalize controllers

	userController := controller.NewUserController(userRepository, storeRepository, userProfileRepository)
	storeController := controller.NewStoreController(storeRepository, productRepository)

	// initialize enforcer

	enforcer := SetupEnforcerPolicies(db)

	// Routes

	apiRoutes := httpRouter.Group("/api")

	{
		apiRoutes.POST("/register/", userController.AddUser(enforcer))
		apiRoutes.POST("/signin/", userController.SignInUser)
	}

	storeRoutes := apiRoutes.Group("/store")
	{
		storeRoutes.GET("/:id/products", storeController.GetProductsByStoreID)
	}

	userProtectedRoutes := apiRoutes.Group("/users", middleware.AuthorizeJWT())
	{
		userProtectedRoutes.GET("/:user/user_profile/", middleware.Authorize("report", "read", enforcer), userController.GetUserProfile)
		//userProtectedRoutes.PUT("/:user/user_profile/", middleware.Authorize("report", "write", enforcer), userController.UpdateUser)
	}

	httpRouter.Run()

}
