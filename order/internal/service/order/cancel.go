package order

import "context"

func (s *orderService) CancelOrderByUUID(ctx context.Context, orderUUID string) (struct{}, error) {
	cancelled, err := s.orderRepository.CancelOrderByUUID(ctx, orderUUID)
	if err != nil {
		return struct{}{}, err
	}
	return cancelled, nil
}
