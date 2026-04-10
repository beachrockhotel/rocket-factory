package service

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.OrderDTO, error)
	PayOrder(ctx context.Context, paymentMethod, orderUUID string) (string, error)
	CancelOrderByUUID(ctx context.Context, orderUUID string) (struct{}, error)
	GetOrderByUUID(ctx context.Context, orderUUID string) (model.OrderDTO, error)
}
