package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/beachrockhotel/rocket-factory/order/internal/converter"
	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	orderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.OrderService.GetOrderByUUID(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.NotFoundError{Message: "order not found"}, nil
		}
		return nil, err
	}

	return &orderV1.GetOrderResponse{
		Order: converter.OrderToProto(order),
	}, nil
}
