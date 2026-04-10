package order

import (
	"context"
)

func (s *orderService) PayOrder(ctx context.Context, paymentMethod, orderUUID string) (string, error) {
	order, err := s.orderRepository.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	transactionUUID, err = s.orderRepository.PayOrder(ctx, transactionUUID, paymentMethod, orderUUID)
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}
