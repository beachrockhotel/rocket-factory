package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUUIDs       = []string{gofakeit.UUID(), gofakeit.UUID()}
		transactionUUID = gofakeit.UUID()
		paymentMethod   = "CARD"
		createdAt       = gofakeit.Date()

		part1 = model.Part{Uuid: partUUIDs[0], Price: gofakeit.Float64()}
		part2 = model.Part{Uuid: partUUIDs[1], Price: gofakeit.Float64()}
		parts = []model.Part{part1, part2}

		expectedTotalPrice = part1.Price + part2.Price

		partsFilter = model.PartsFilter{
			Uuids: partUUIDs,
		}

		orderRequest = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUUIDs: partUUIDs,
		}

		expectedRepoRequest = model.CreateOrderParams{
			UserUUID:   userUUID,
			PartUUIDs:  partUUIDs,
			TotalPrice: expectedTotalPrice,
			Status:     model.OrderStatusPendingPayment,
		}

		expectedOrderResponse = model.OrderDTO{
			UUID:            orderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUUIDs,
			TotalPrice:      expectedTotalPrice,
			TransactionUUID: &transactionUUID,
			PaymentMethod:   &paymentMethod,
			Status:          model.OrderStatusPendingPayment,
			CreatedAt:       createdAt,
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsFilter).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, expectedRepoRequest).Return(expectedOrderResponse, nil)

	res, err := s.service.CreateOrder(s.ctx, orderRequest)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
	s.Require().Equal(expectedOrderResponse, res)
}

func (s *ServiceSuite) TestCreateOrderInventoryClientError() {
	var (
		userUUID       = gofakeit.UUID()
		partUUIDs      = []string{gofakeit.UUID(), gofakeit.UUID()}
		inventoryError = gofakeit.Error()

		partsFilter = model.PartsFilter{
			Uuids: partUUIDs,
		}

		orderRequest = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUUIDs: partUUIDs,
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsFilter).Return(nil, inventoryError)
	res, err := s.service.CreateOrder(s.ctx, orderRequest)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "inventory service error")
	s.Require().ErrorContains(err, inventoryError.Error())
	s.Require().Empty(res)
}

func (s *ServiceSuite) TestCreateOrderError() {
	var (
		repoErr   = gofakeit.Error()
		userUUID  = gofakeit.UUID()
		partUUIDs = []string{gofakeit.UUID(), gofakeit.UUID()}
		part1     = model.Part{Uuid: partUUIDs[0], Price: gofakeit.Float64()}
		part2     = model.Part{Uuid: partUUIDs[1], Price: gofakeit.Float64()}
		parts     = []model.Part{part1, part2}

		orderRequest = model.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUUIDs: partUUIDs,
		}

		listRequest = model.PartsFilter{
			Uuids: partUUIDs,
		}

		expectedTotalPrice = part1.Price + part2.Price

		expectedRepoRequest = model.CreateOrderParams{
			UserUUID:   userUUID,
			PartUUIDs:  partUUIDs,
			TotalPrice: expectedTotalPrice,
			Status:     model.OrderStatusPendingPayment,
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, listRequest).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, expectedRepoRequest).Return(model.OrderDTO{}, repoErr)

	resp, err := s.service.CreateOrder(s.ctx, orderRequest)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(resp)
}
