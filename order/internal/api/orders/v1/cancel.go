package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	orderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	_, err := a.OrderService.CancelOrderByUUID(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.NotFoundError{Message: "order not found"}, nil
		}
	}
	return &orderV1.CancelOrderByUUIDNoContent{}, nil
}
