package v1

import (
	"github.com/beachrockhotel/rocket-factory/payment/internal/service"
	paymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	payment service.PaymentService
}

func NewApi(payment service.PaymentService) *api {
	return &api{
		payment: payment,
	}
}
