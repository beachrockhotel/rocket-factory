package inventory

import (
	"context"
	"strings"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
	"github.com/beachrockhotel/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/beachrockhotel/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	result := make([]repoModel.Part, 0, len(r.parts))

	uuidSet := make(map[string]struct{}, len(filter.Uuids))
	for _, partUUID := range filter.Uuids {
		uuidSet[partUUID] = struct{}{}
	}

	nameSet := make(map[string]struct{}, len(filter.Names))
	for _, name := range filter.Names {
		nameSet[name] = struct{}{}
	}

	categorySet := make(map[int32]struct{}, len(filter.Categories))
	for _, category := range filter.Categories {
		categorySet[category] = struct{}{}
	}

	manufacturerCountriesSet := make(map[string]struct{}, len(filter.ManufacturerCountries))
	for _, manufacturerCountry := range filter.ManufacturerCountries {
		manufacturerCountriesSet[manufacturerCountry] = struct{}{}
	}

	tagsSet := make(map[string]struct{}, len(filter.Tags))
	for _, tag := range filter.Tags {
		tagsSet[tag] = struct{}{}
	}

	for _, part := range r.parts {
		if len(uuidSet) > 0 {
			if _, ok := uuidSet[part.Uuid]; !ok {
				continue
			}
		}

		if len(nameSet) > 0 {
			if _, ok := nameSet[part.Name]; !ok {
				continue
			}
		}

		if len(categorySet) > 0 {
			if _, ok := categorySet[part.Category]; !ok {
				continue
			}
		}

		if len(manufacturerCountriesSet) > 0 {
			if _, ok := manufacturerCountriesSet[part.Manufacturer.Country]; !ok {
				continue
			}
		}

		if len(tagsSet) > 0 {
			foundTag := false
			for _, tag := range part.Tags {
				if _, ok := tagsSet[strings.ToLower(tag)]; ok {
					foundTag = true
					break
				}
			}
			if !foundTag {
				continue
			}
		}
		result = append(result, part)
	}
	return converter.ListPartsToModel(result), nil
}
