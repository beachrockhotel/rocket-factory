package v1

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/payment/internal/model"
	paymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *APISuite) TestPayOrderSuccess() {
	var (
		expectedUUID  = gofakeit.UUID()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD

		payOrderRequest = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		}

		expectedModelRequest = model.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: "PAYMENT_METHOD_CARD",
		}
	)

	a.paymentService.
		On("PayOrder", a.ctx, expectedModelRequest).
		Return(expectedUUID, nil)

	res, err := a.api.PayOrder(a.ctx, payOrderRequest)
	a.Require().NoError(err)
	a.Require().NotNil(res)
	a.Require().Equal(expectedUUID, res.TransactionUuid)
}

func (a *APISuite) TestPayOrderError() {
	var (
		serviceErr    = gofakeit.Error()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD

		payOrderRequest = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		}

		exceptedModelRequest = model.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: "PAYMENT_METHOD_CARD",
		}
	)

	a.paymentService.On("PayOrder", a.ctx, exceptedModelRequest).Return("", serviceErr)

	uuid, err := a.api.PayOrder(a.ctx, payOrderRequest)
	a.Require().Error(err)
	a.Require().ErrorIs(err, serviceErr)
	a.Require().Empty(uuid)
}
