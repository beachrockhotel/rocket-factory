package converter

import (
	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	repoModel "github.com/beachrockhotel/rocket-factory/order/internal/repository/model"
)

func OrderToRepoModel(order model.OrderDTO) repoModel.OrderDTO {
	return repoModel.OrderDTO{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func OrderFromRepoToModel(order repoModel.OrderDTO) model.OrderDTO {
	return model.OrderDTO{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
	}
}
