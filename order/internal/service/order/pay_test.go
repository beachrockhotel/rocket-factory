package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		paymentMethod   = "CARD"
		uuid            = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUuids       = []string{gofakeit.UUID(), gofakeit.UUID()}
		totalPrice      = gofakeit.Float64()
		transactionUUID = gofakeit.UUID()
		status          = "PAID"
		createdAt       = gofakeit.Date()
		updatedAt       = gofakeit.Date()
		getOrderResp    = model.OrderDTO{
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
	s.orderRepository.On("GetOrderByUUID", s.ctx, uuid).Return(getOrderResp, nil)
	s.paymentClient.On("PayOrder", s.ctx, uuid, getOrderResp.UserUUID, paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("PayOrder", s.ctx, transactionUUID, paymentMethod, uuid).Return(transactionUUID, nil)

	resp, err := s.service.PayOrder(s.ctx, paymentMethod, uuid)

	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().Equal(transactionUUID, resp)
}

func (s *ServiceSuite) TestPayOrderOrderRepoError() {
	var (
		repoError     = gofakeit.Error()
		paymentMethod = "CARD"
		uuid          = gofakeit.UUID()
	)
	s.orderRepository.On("GetOrderByUUID", s.ctx, uuid).Return(model.OrderDTO{}, repoError)
	resp, err := s.service.PayOrder(s.ctx, paymentMethod, uuid)

	s.Require().Error(err)
	s.Require().Empty(resp)
	s.Require().Equal(err, repoError)
}

func (s *ServiceSuite) TestPayOrderPaymentClientError() {
	var (
		repoError       = gofakeit.Error()
		paymentMethod   = "CARD"
		uuid            = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUuids       = []string{gofakeit.UUID(), gofakeit.UUID()}
		totalPrice      = gofakeit.Float64()
		transactionUUID = gofakeit.UUID()
		status          = "PAID"
		createdAt       = gofakeit.Date()
		updatedAt       = gofakeit.Date()
		getOrderResp    = model.OrderDTO{
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
	s.orderRepository.On("GetOrderByUUID", s.ctx, uuid).Return(getOrderResp, nil)
	s.paymentClient.On("PayOrder", s.ctx, uuid, getOrderResp.UserUUID, paymentMethod).Return("", repoError)

	resp, err := s.service.PayOrder(s.ctx, paymentMethod, uuid)

	s.Require().Error(err)
	s.Require().Empty(resp)
	s.Require().Equal(err, repoError)
}

func (s *ServiceSuite) TestPayOrderRepositoryError() {
	var (
		repoError       = gofakeit.Error()
		paymentMethod   = "CARD"
		uuid            = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUuids       = []string{gofakeit.UUID(), gofakeit.UUID()}
		totalPrice      = gofakeit.Float64()
		transactionUUID = gofakeit.UUID()
		status          = "PAID"
		createdAt       = gofakeit.Date()
		updatedAt       = gofakeit.Date()
		getOrderResp    = model.OrderDTO{
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
	s.orderRepository.On("GetOrderByUUID", s.ctx, uuid).Return(getOrderResp, nil)
	s.paymentClient.On("PayOrder", s.ctx, uuid, getOrderResp.UserUUID, paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("PayOrder", s.ctx, transactionUUID, paymentMethod, uuid).Return("", repoError)

	resp, err := s.service.PayOrder(s.ctx, paymentMethod, uuid)

	s.Require().Error(err)
	s.Require().Empty(resp)
	s.Require().Equal(err, repoError)
}
