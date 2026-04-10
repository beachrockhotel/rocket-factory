package converter

import (
	"github.com/beachrockhotel/rocket-factory/payment/internal/model"
	repoModel "github.com/beachrockhotel/rocket-factory/payment/internal/repository/model"
)

func PayOrderToRepo(req model.PayOrderRequest) repoModel.PayOrderRequest {
	return repoModel.PayOrderRequest{
		OrderUuid:     req.OrderUuid,
		UserUUID:      req.UserUUID,
		PaymentMethod: req.PaymentMethod,
	}
}
