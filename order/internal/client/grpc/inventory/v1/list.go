package v1

import (
	"context"

	clientConverter "github.com/beachrockhotel/rocket-factory/order/internal/client/converter"
	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	listPartsResponse, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.FilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	parts := listPartsResponse.Parts

	return clientConverter.ArrayPartToModel(parts), nil
}
