package order

import (
	"github.com/beachrockhotel/rocket-factory/order/internal/client/grpc"
	"github.com/beachrockhotel/rocket-factory/order/internal/repository"
)

type orderService struct {
	orderRepository repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *orderService {
	return &orderService{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
