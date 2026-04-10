package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/beachrockhotel/rocket-factory/order/internal/model"
	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartToModel(part *inventoryV1.Part) model.Part {
	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		t := part.UpdatedAt.AsTime()
		updatedAt = &t
	}
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimension:     DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func DimensionsToModel(dimensions *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
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

func ArrayPartToModel(parts []*inventoryV1.Part) []model.Part {
	data := make([]model.Part, 0, len(parts))
	for _, v := range parts {
		data = append(data, PartToModel(v))
	}
	return data
}

func FilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	uuidSet := make([]string, 0, len(filter.Uuids))
	if len(filter.Uuids) > 0 {
		uuidSet = append(uuidSet, filter.Uuids...)
	}
	nameSet := make([]string, 0, len(filter.Names))
	if len(filter.Names) > 0 {
		uuidSet = append(uuidSet, filter.Names...)
	}
	categorySet := make([]inventoryV1.Category, 0, len(filter.Categories))
	if len(filter.Categories) > 0 {
		for _, v := range filter.Categories {
			categorySet = append(categorySet, CategoryToProto(v))
		}
	}
	manufacturerSet := make([]string, 0, len(filter.ManufacturerCountries))
	if len(filter.ManufacturerCountries) > 0 {
		manufacturerSet = append(manufacturerSet, filter.ManufacturerCountries...)
	}

	tagSet := make([]string, 0, len(filter.Tags))
	if len(filter.Tags) > 0 {
		tagSet = append(tagSet, filter.Tags...)
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuidSet,
		Names:                 nameSet,
		Categories:            categorySet,
		ManufacturerCountries: manufacturerSet,
		Tags:                  tagSet,
	}
}

func CategoryToProto(category int32) inventoryV1.Category {
	switch category {
	case 0:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	case 1:
		return inventoryV1.Category_CATEGORY_ENGINE
	case 2:
		return inventoryV1.Category_CATEGORY_FUEL
	case 3:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case 4:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return 0
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

func InventoryToProto(part model.Part) *inventoryV1.Part {
	createdAt := timestamppb.New(part.CreatedAt)
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      safeCategory(part.Category),
		Dimensions:    DimensionsToProto(part.Dimension),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToProto(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func safeCategory(category int32) inventoryV1.Category {
	switch inventoryV1.Category(category) {
	case inventoryV1.Category_CATEGORY_UNSPECIFIED,
		inventoryV1.Category_CATEGORY_ENGINE,
		inventoryV1.Category_CATEGORY_FUEL,
		inventoryV1.Category_CATEGORY_PORTHOLE,
		inventoryV1.Category_CATEGORY_WING:
		return inventoryV1.Category(category)
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func DimensionsToProto(dimensions model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToProto(manufacturer model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToProto(metadata map[string]any) map[string]*inventoryV1.Value {
	res := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		switch v := value.(type) {
		case string:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{StringValue: v},
			}
		case int64, int32, int:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_Int64Value{Int64Value: toInt64(v)},
			}
		case float32:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: float64(v)},
			}
		case float64:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: v},
			}
		case bool:
			res[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_BoolValue{BoolValue: v},
			}
		default:

		}
	}

	return res
}

func toInt64(value any) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}
