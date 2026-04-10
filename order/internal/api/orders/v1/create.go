package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/beachrockhotel/rocket-factory/order/internal/converter"
	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	orderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	creatOrderResp, err := a.OrderService.CreateOrder(ctx, converter.OrderToModel(req))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return &orderV1.ServiceUnavailableError{Message: "order service timeout"}, nil
		}
		if errors.Is(err, model.ErrBadRequest) {
			return &orderV1.BadRequestError{Message: "some parts not found"}, nil
		}
		return &orderV1.InternalServerError{
			Message: fmt.Sprintf("order service errors %s", err.Error()),
		}, nil
	}
	return &orderV1.CreateOrderResponse{
		Order: converter.OrderToProto(creatOrderResp),
	}, nil
}
