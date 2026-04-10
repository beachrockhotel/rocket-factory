package converter

import (
	"github.com/google/uuid"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	orderV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderToModel(req *orderV1.CreateOrderRequest) model.CreateOrderRequest {
	data := make([]string, 0, len(req.PartUuids))
	userUUID := req.UserUUID.String()

	for _, part := range req.PartUuids {
		data = append(data, part.String())
	}

	return model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUUIDs: data,
	}
}

func OrderToProto(order model.OrderDTO) orderV1.OrderDto {
	partUUIDs := make([]uuid.UUID, 0, len(order.PartUuids))
	for _, part := range order.PartUuids {
		partUUIDs = append(partUUIDs, uuid.MustParse(part))
	}

	res := orderV1.OrderDto{
		UUID:       uuid.MustParse(order.UUID),
		UserUUID:   uuid.MustParse(order.UserUUID),
		PartUuids:  partUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     OrderStatus(order.Status),
		CreatedAt:  order.CreatedAt,
	}

	if order.TransactionUUID != nil {
		res.TransactionUUID = orderV1.NewOptNilUUID(uuid.MustParse(*order.TransactionUUID))
	}

	if order.PaymentMethod != nil {
		res.PaymentMethod = orderV1.NewOptPaymentMethod(PaymentMethodToProto(*order.PaymentMethod))
	}

	if order.UpdatedAt != nil {
		res.UpdatedAt = orderV1.NewOptDateTime(*order.UpdatedAt)
	}

	return res
}

func PaymentMethodToProto(paymentMethod string) orderV1.PaymentMethod {
	switch paymentMethod {
	case "PAYMENT_METHOD_UNSPECIFIED":
		return orderV1.PaymentMethodPAYMENTMETHODUNSPECIFIED
	case "PAYMENT_METHOD_CARD":
		return orderV1.PaymentMethodPAYMENTMETHODCARD
	case "PAYMENT_METHOD_SBP":
		return orderV1.PaymentMethodPAYMENTMETHODSBP
	case "PAYMENT_METHOD_CREDIT_CARD":
		return orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	case "PAYMENT_METHOD_INVESTOR_MONEY":
		return orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		return orderV1.PaymentMethodPAYMENTMETHODUNSPECIFIED
	}
}

func PaymentMethodToModel(paymentMethod orderV1.PaymentMethod) string {
	switch paymentMethod {
	case orderV1.PaymentMethodPAYMENTMETHODUNSPECIFIED:
		return "PAYMENT_METHOD_UNSPECIFIED"
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return "PAYMENT_METHOD_CARD"
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return "PAYMENT_METHOD_SBP"
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return "PAYMENT_METHOD_CREDIT_CARD"
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return "PAYMENT_METHOD_INVESTOR_MONEY"
	default:
		return "PAYMENT_METHOD_UNSPECIFIED"
	}
}

func OrderStatus(status string) orderV1.OrderStatus {
	switch status {
	case "PENDING_PAYMENT":
		return orderV1.OrderStatusPENDINGPAYMENT
	case "PAID":
		return orderV1.OrderStatusPAID
	case "CANCELLED":
		return orderV1.OrderStatusCANCELLED
	default:
		return ""
	}
}
