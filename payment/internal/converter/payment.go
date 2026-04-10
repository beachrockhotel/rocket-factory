package converter

import (
	"github.com/beachrockhotel/rocket-factory/payment/internal/model"
	paymentV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/payment/v1"
)

func ReqToModel(req *paymentV1.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		OrderUuid:     req.OrderUuid,
		UserUUID:      req.UserUuid,
		PaymentMethod: PaymentMethodToModel(req.PaymentMethod),
	}
}

func PaymentMethodToModel(pm paymentV1.PaymentMethod) string {
	switch pm {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		return "PAYMENT_METHOD_UNSPECIFIED"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return "PAYMENT_METHOD_CARD"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return "PAYMENT_METHOD_SBP"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return "PAYMENT_METHOD_CREDIT_CARD"
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return "PAYMENT_METHOD_INVESTOR_MONEY"
	default:
		return "PAYMENT_METHOD_UNSPECIFIED"
	}
}
