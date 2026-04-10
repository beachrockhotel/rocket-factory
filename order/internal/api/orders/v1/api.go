package v1

import "github.com/beachrockhotel/rocket-factory/order/internal/service"

type OrderAPI struct {
	OrderService service.OrderService
}

func NewOrderAPI(orderService service.OrderService) *OrderAPI {
	return &OrderAPI{
		OrderService: orderService,
	}
}
