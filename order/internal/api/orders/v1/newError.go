package v1

import (
	"context"
	"net/http"

	orderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.GenericErrorCode520,
			Message: err.Error(),
		},
	}
}
