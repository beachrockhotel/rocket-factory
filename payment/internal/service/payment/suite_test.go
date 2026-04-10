package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/beachrockhotel/rocket-factory/payment/internal/repository/mocks"
)

//nolint:containedctx
type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	paymentRepository *mocks.PaymentRepository

	service *paymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentRepository = mocks.NewPaymentRepository(s.T())
	s.service = NewPaymentService(
		s.paymentRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
