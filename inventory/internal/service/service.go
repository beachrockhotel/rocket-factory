package service

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
)

type InventoryService interface {
	GetPart(_ context.Context, uuid string) (model.Part, error)
	ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error)
}
