package v1

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/payment/internal/converter"
	paymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID, err := a.payment.PayOrder(ctx, converter.ReqToModel(req))
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
