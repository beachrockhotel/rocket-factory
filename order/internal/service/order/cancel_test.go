package order

import "github.com/brianvoe/gofakeit/v7"

func (s *ServiceSuite) TestCancelOrderByUUIDSuccess() {
	orderUUID := gofakeit.UUID()
	s.orderRepository.On("CancelOrderByUUID", s.ctx, orderUUID).Return(struct{}{}, nil)

	res, err := s.service.CancelOrderByUUID(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(struct{}{}, res)
}

func (s *ServiceSuite) TestCancelOrderByUUIDError() {
	var (
		orderUUID = gofakeit.UUID()
		repoError = gofakeit.Error()
	)
	s.orderRepository.On("CancelOrderByUUID", s.ctx, orderUUID).Return(struct{}{}, repoError)

	res, err := s.service.CancelOrderByUUID(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().Equal(repoError, err)
}
