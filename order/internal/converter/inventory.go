package converter

import (
	"time"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartToModel(part *inventoryV1.Part) model.Part {
	var updatedAt time.Time
	if part.UpdatedAt != nil {
		updatedAt = part.UpdatedAt.AsTime()
	}

	createdAt := part.CreatedAt.AsTime()

	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimension:     DimensionToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     &updatedAt,
	}
}

func CategoryToModel(category inventoryV1.Category) int32 {
	switch category {
	case inventoryV1.Category_CATEGORY_UNSPECIFIED:
		return 0
	case inventoryV1.Category_CATEGORY_ENGINE:
		return 1
	case inventoryV1.Category_CATEGORY_FUEL:
		return 2
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return 3
	case inventoryV1.Category_CATEGORY_WING:
		return 4
	default:
		return 0
	}
}

func DimensionToModel(dimension *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataValueToModel(metadata *inventoryV1.Value) any {
	if metadata == nil {
		return nil
	}

	switch v := metadata.GetKind().(type) {
	case *inventoryV1.Value_StringValue:
		return v.StringValue
	case *inventoryV1.Value_DoubleValue:
		return v.DoubleValue
	case *inventoryV1.Value_BoolValue:
		return v.BoolValue
	case *inventoryV1.Value_Int64Value:
		return v.Int64Value
	default:
		return nil
	}
}

func PartMetadataToModel(metadata map[string]*inventoryV1.Value) map[string]any {
	data := make(map[string]any)
	for k, v := range metadata {
		data[k] = MetadataValueToModel(v)
	}
	return data
}
