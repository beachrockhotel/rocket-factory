package repository

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.CreateOrderParams) (model.OrderDTO, error)
	PayOrder(ctx context.Context, transactionUUID, paymentMethod, orderUUID string) (string, error)
	CancelOrderByUUID(ctx context.Context, orderUUID string) (struct{}, error)
	GetOrderByUUID(ctx context.Context, orderUUID string) (model.OrderDTO, error)
}
