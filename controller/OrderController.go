package controller

type OrderController interface {
	CreateOrder()
	TakeOrder()
	RejectOrder()
	DoneOrder()
	GetOrders()
	GetOrder()
}

type orderConrtoller struct {
}
