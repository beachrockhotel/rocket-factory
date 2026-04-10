package order

import (
	"context"
	"log"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (r *orderRepository) CancelOrderByUUID(_ context.Context, orderUUID string) (struct{}, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	order, ok := r.data[orderUUID]
	if !ok {
		return struct{}{}, model.ErrNotFound
	}

	if order.Status != model.OrderStatusPendingPayment {
		return struct{}{}, model.ErrConflict
	}
	order.Status = model.OrderStatusCancelled
	r.data[orderUUID] = order
	log.Printf(`
	💳 [Order Canceled]
	• 🆔 Order UUID: %s
	• 👤 User UUID: %s
	• 💰 Status: %s
	`, order.UUID, order.UserUUID, order.Status,
	)
	return struct{}{}, nil
}
