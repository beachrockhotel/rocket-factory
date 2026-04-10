package v1

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/beachrockhotel/rocket-factory/order/internal/converter"
	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	orderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transactionUUID, err := a.OrderService.PayOrder(ctx, converter.PaymentMethodToModel(req.PaymentMethod), params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &orderV1.NotFoundError{Message: "order not found"}, nil
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return &orderV1.ServiceUnavailableError{
				Message: "repository service timeout",
			}, nil
		}

		return &orderV1.InternalServerError{
			Message: fmt.Sprintf("repository service error: %s", err.Error()),
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID),
	}, nil
}
