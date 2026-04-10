package order

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (s *orderService) GetOrderByUUID(ctx context.Context, uuid string) (model.OrderDTO, error) {
	order, err := s.orderRepository.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return model.OrderDTO{}, err
	}

	return order, nil
}
