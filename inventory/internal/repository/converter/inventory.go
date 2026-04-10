package converter

import (
	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
	repoModel "github.com/beachrockhotel/rocket-factory/inventory/internal/repository/model"
)

func InventoryToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Dimensions:    DimensionToRepoModel(part.Dimension),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     *part.UpdatedAt,
	}
}

func DimensionToRepoModel(dimension model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func ManufacturerToRepoModel(manufacturer model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func InventoryToModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Dimension:     DimensionToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     &part.UpdatedAt,
	}
}

func DimensionToModel(dimension repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func ManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ListPartsToModel(partList []repoModel.Part) []model.Part {
	res := make([]model.Part, 0, len(partList))
	for _, part := range partList {
		convertedModel := InventoryToModel(part)
		res = append(res, convertedModel)
	}
	return res
}
