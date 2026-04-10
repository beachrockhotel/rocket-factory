package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	mocks "github.com/beachrockhotel/rocket-factory/payment/internal/service/mocks"
)

//nolint:containedctx
type APISuite struct {
	suite.Suite

	ctx context.Context

	paymentService *mocks.PaymentService

	api *api
}

func (a *APISuite) SetupTest() {
	a.ctx = context.Background()

	a.paymentService = mocks.NewPaymentService(a.T())

	a.api = NewApi(
		a.paymentService,
	)
}

func (a *APISuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
