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

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//add policy for client and story_admin
	if hasPolicy := enforcer.HasPolicy("CLIENT", "profile", "read"); !hasPolicy {
		enforcer.AddPolicy("CLIENT", "profile", "read")
	}
	if hasPolicy := enforcer.HasPolicy("CLIENT", "profile", "update"); !hasPolicy {
		enforcer.AddPolicy("CLIENT", "profile", "update")
	}
	if hasPolicy := enforcer.HasPolicy("STOREADMIN", "profile", "read"); !hasPolicy {
		enforcer.AddPolicy("CLIENT", "profile", "read")
	}
	if hasPolicy := enforcer.HasPolicy("CLIENT", "profile", "update"); !hasPolicy {
		enforcer.AddPolicy("CLIENT", "profile", "update")
	}

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
		apiRoutes.POST(
			"/register/",
			userController.AddUser(enforcer),
		)
		apiRoutes.POST(
			"/signin/",
			userController.SignInUser,
		)
	}

	storeRoutes := apiRoutes.Group("/store")
	{
		storeRoutes.GET("/:id/products", storeController.GetProductsByStoreID)
	}

	userProtectedRoutes := apiRoutes.Group("/users", middleware.AuthorizeJWT())
	{
		userProtectedRoutes.GET(
			"/my/user_profile/",
			middleware.Authorize("profile", "read", enforcer),
			userController.GetUserProfile,
		)
		userProtectedRoutes.PUT(
			"/my/user_profile/",
			middleware.Authorize("profile", "update", enforcer),
			userController.UpdateUserProfile,
		)
	}

	httpRouter.Run()

}
