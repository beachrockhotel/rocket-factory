package inventory

import (
	"context"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
