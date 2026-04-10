package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/repository/mocks"
)

//nolint:containedctx
type ServiceSuite struct {
	suite.Suite
	ctx        context.Context
	repository *mocks.InventoryRepository
	service    *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repository = mocks.NewInventoryRepository(s.T())
	s.service = NewService(s.repository)
}

func (s *ServiceSuite) TearDownTest() {}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
