package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/beachrockhotel/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/beachrockhotel/rocket-factory/order/internal/repository/mocks"
)

//nolint:containedctx
type ServiceSuite struct {
	suite.Suite
	ctx             context.Context
	orderRepository *mocks.OrderRepository
	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient
	service         *orderService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
