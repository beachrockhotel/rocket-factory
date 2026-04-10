package v1

import (
	"context"

	generatedPaymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: generatedPaymentV1.PaymentMethod(generatedPaymentV1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}
