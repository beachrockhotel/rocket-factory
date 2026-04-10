package inventory

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
	"github.com/beachrockhotel/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	modelPart := converter.InventoryToModel(part)

	return modelPart, nil
}
