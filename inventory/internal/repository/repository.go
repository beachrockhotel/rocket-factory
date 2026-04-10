package repository

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(_ context.Context, uuid string) (model.Part, error)
	ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error)
}
