package order

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	"github.com/beachrockhotel/rocket-factory/order/internal/repository/converter"
)

func (r *orderRepository) GetOrderByUUID(_ context.Context, uuid string) (model.OrderDTO, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	order, ok := r.data[uuid]
	if !ok {
		return model.OrderDTO{}, model.ErrNotFound
	}

	return converter.OrderFromRepoToModel(order), nil
}
