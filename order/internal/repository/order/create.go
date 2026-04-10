package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	repoConverter "github.com/beachrockhotel/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/beachrockhotel/rocket-factory/order/internal/repository/model"
)

func (r *orderRepository) CreateOrder(_ context.Context, order model.CreateOrderParams) (model.OrderDTO, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	newUUID := uuid.NewString()

	newOrder := repoModel.OrderDTO{
		UUID:       newUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  time.Now(),
	}

	r.data[newUUID] = newOrder

	return repoConverter.OrderFromRepoToModel(r.data[newUUID]), nil
}
