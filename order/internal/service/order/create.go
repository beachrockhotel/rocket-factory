package order

import (
	"context"
	"log"

	"github.com/go-faster/errors"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (s *orderService) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (model.OrderDTO, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: req.PartUUIDs,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return model.OrderDTO{}, errors.New("inventory service timeout")
		}

		return model.OrderDTO{}, errors.Errorf("inventory service error: %s", err.Error())
	}

	if len(parts) != len(req.PartUUIDs) {
		return model.OrderDTO{}, model.ErrNotFound
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	createParams := model.CreateOrderParams{
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}

	order, err := s.orderRepository.CreateOrder(ctx, createParams)
	if err != nil {
		return model.OrderDTO{}, err
	}

	log.Printf(`
	💳 [Order Created]
	• 🆔 Order UUID: %s
	• 👤 User UUID: %s
	• 💰 Part UUIDs: %v
	• 💰 Total Price: %f
	• 💰 Status: %s
	• 💰 CreatedAt: %v
	`, order.UUID, order.UserUUID, order.PartUuids, order.TotalPrice, order.Status, order.CreatedAt)

	return order, nil
}
