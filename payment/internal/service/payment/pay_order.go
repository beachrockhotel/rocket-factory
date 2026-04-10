package payment

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/payment/internal/model"
)

func (s *paymentService) PayOrder(ctx context.Context, request model.PayOrderRequest) (string, error) {
	payOrder, err := s.repo.PayOrder(ctx, request)
	if err != nil {
		return "", err
	}
	return payOrder, nil
}
