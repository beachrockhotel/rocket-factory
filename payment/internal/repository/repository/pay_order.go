package repository

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/beachrockhotel/rocket-factory/payment/internal/model"
	"github.com/beachrockhotel/rocket-factory/payment/internal/repository/converter"
)

func (r *paymentRepository) PayOrder(_ context.Context, req model.PayOrderRequest) (string, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	orderRepo := converter.PayOrderToRepo(req)
	log.Printf(`
	💳 [Order Paid]
	• 🆔 Order UUID: %s
	• 👤 User UUID: %s
	• 💰 Payment Method: %s`,
		orderRepo.OrderUuid, orderRepo.UserUUID, orderRepo.PaymentMethod)

	transactionUUID := uuid.NewString()
	log.Print(transactionUUID)

	return transactionUUID, nil
}
