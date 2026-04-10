package order

import (
	"context"
	"log"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (r *orderRepository) PayOrder(_ context.Context, transactionUUID, paymentMethod, orderUUID string) (string, error) {
	r.mtx.Lock()
	order, ok := r.data[orderUUID]
	if !ok {
		return "", model.ErrNotFound
	}
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethod
	order.Status = model.OrderStatusPaid
	r.data[orderUUID] = order
	r.mtx.Unlock()

	log.Printf(`
	💳 [Order Paid]
	• 🆔 Order UUID: %s
	• 👤 User UUID: %s
	• 💰 Transaction UUID: %s
	• 💰 Payment Method: %s
	• 💰 Status: %s
	`, order.UUID, order.UserUUID, *order.TransactionUUID, *order.PaymentMethod, order.Status,
	)

	return transactionUUID, nil
}
