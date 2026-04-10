package payment

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		transactionUUID = gofakeit.UUID()
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		paymentMethod   = "CARD"

		payment = model.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: paymentMethod,
		}
	)

	s.paymentRepository.On("PayOrder", s.ctx, payment).Return(transactionUUID, nil)

	res, err := s.service.PayOrder(s.ctx, payment)
	s.NoError(err)
	s.Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderError() {
	var (
		repoErr       = gofakeit.Error()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = "CARD"

		payment = model.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: paymentMethod,
		}
	)
	s.paymentRepository.On("PayOrder", s.ctx, payment).Return("", repoErr)

	uuid, err := s.service.PayOrder(s.ctx, payment)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(uuid)
}
