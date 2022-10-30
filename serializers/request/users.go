package request

type CreateUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
	StoreID  int    `json:"store_id"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserProfileInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	City      string `json:"city" binding:"required"`
	Address   string `json:"address" binding:"required"`
}
