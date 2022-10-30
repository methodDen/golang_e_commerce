package response

type UserProfileResponse struct {
	ID        uint   `json:"id" binding:"required"`
	UserID    uint   `json:"user_id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	City      string `json:"city" binding:"required"`
	Address   string `json:"address" binding:"required"`
}
