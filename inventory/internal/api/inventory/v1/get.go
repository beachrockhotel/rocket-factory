package v1

import (
	"context"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/converter"
	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *inventoryAPI) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part not found")
		}
	}
	return &inventoryV1.GetPartResponse{
		Part: converter.InventoryToProto(part),
	}, nil
}
