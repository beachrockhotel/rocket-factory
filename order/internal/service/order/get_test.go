package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderByUUIDSuccess() {
	var (
		uuid            = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUuids       = []string{gofakeit.UUID(), gofakeit.UUID()}
		totalPrice      = gofakeit.Float64()
		transactionUUID = gofakeit.UUID()
		paymentMethod   = "CARD"
		status          = "PENDING_PAYMENT"
		createdAt       = gofakeit.Date()
		updatedAt       = gofakeit.Date()

		exceptedOrder = model.OrderDTO{
			UUID:            uuid,
			UserUUID:        userUUID,
			PartUuids:       partUuids,
			TotalPrice:      totalPrice,
			TransactionUUID: &transactionUUID,
			PaymentMethod:   &paymentMethod,
			Status:          status,
			CreatedAt:       createdAt,
			UpdatedAt:       &updatedAt,
		}
	)

	s.orderRepository.On("GetOrderByUUID", s.ctx, uuid).Return(exceptedOrder, nil)

	order, err := s.service.GetOrderByUUID(s.ctx, uuid)
	s.Require().NoError(err)
	s.Require().NotEmpty(order)
	s.Require().Equal(exceptedOrder, order)
}

func (s *ServiceSuite) TestCreateOrderByUUIDError() {
	var (
		errorRepo = gofakeit.Error()
		uuid      = gofakeit.UUID()
	)
	s.orderRepository.On("GetOrderByUUID", s.ctx, uuid).Return(model.OrderDTO{}, errorRepo)

	order, err := s.service.GetOrderByUUID(s.ctx, uuid)

	s.Require().Error(err)
	s.Require().Equal(model.OrderDTO{}, order)
	s.Require().Equal(errorRepo, err)
}
