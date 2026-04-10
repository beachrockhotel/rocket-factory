package v1

import (
	"github.com/beachrockhotel/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
)

type inventoryAPI struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *inventoryAPI {
	return &inventoryAPI{
		inventoryService: inventoryService,
	}
}
