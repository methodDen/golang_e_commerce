package request

type CreateProductInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Categories  []int  `json:"categories" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Status      string `json:"status" binding:"required"`
}
