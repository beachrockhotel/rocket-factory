package v1

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *inventoryAPI) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.FilterToModel(req.Filter))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.ListPartsToProto(parts),
	}, nil
}
