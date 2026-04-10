package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		uuid1         = gofakeit.UUID()
		uuid2         = gofakeit.UUID()
		name1         = gofakeit.Name()
		name2         = gofakeit.Name()
		description   = gofakeit.ProductDescription()
		price         = gofakeit.Float64()
		stockQuantity = gofakeit.Int64()
		category1     = gofakeit.Int32()
		category2     = gofakeit.Int32()
		dimension     = model.Dimensions{
			Length: gofakeit.Float64(),
			Width:  gofakeit.Float64(),
			Height: gofakeit.Float64(),
			Weight: gofakeit.Float64(),
		}
		manufacturer1 = model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.AppName(),
		}
		manufacturer2 = model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.AppName(),
		}
		tag1     = []string{gofakeit.Gamertag(), gofakeit.Gamertag()}
		tag2     = []string{gofakeit.Gamertag(), gofakeit.Gamertag()}
		metadata = map[string]any{
			gofakeit.Name(): gofakeit.Country(),
			gofakeit.Name(): gofakeit.Country(),
			gofakeit.Name(): gofakeit.Country(),
		}
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()
		part1     = model.Part{
			Uuid:          uuid1,
			Name:          name1,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      category1,
			Dimension:     dimension,
			Manufacturer:  manufacturer1,
			Tags:          tag1,
			Metadata:      metadata,
			CreatedAt:     createdAt,
			UpdatedAt:     &updatedAt,
		}

		part2 = model.Part{
			Uuid:          uuid2,
			Name:          name2,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      category2,
			Dimension:     dimension,
			Manufacturer:  manufacturer2,
			Tags:          tag2,
			Metadata:      metadata,
			CreatedAt:     createdAt,
			UpdatedAt:     &updatedAt,
		}
		uuids                 = []string{uuid1, uuid2}
		names                 = []string{name1, name2}
		categories            = []int32{category1, category2}
		manufacturerCountries = []string{manufacturer1.Country, manufacturer2.Country}
		tags                  = []string{tag1[0], tag2[0]}

		partsFilter = model.PartsFilter{
			Uuids:                 uuids,
			Names:                 names,
			Categories:            categories,
			ManufacturerCountries: manufacturerCountries,
			Tags:                  tags,
		}

		parts = []model.Part{part1, part2}
	)
	s.repository.On("ListParts", s.ctx, partsFilter).Return(parts, nil)

	resp, err := s.service.ListParts(s.ctx, partsFilter)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().Equal(parts, resp)
}

func (s *ServiceSuite) TestListPartError() {
	var (
		uuid1         = gofakeit.UUID()
		uuid2         = gofakeit.UUID()
		name1         = gofakeit.Name()
		name2         = gofakeit.Name()
		category1     = gofakeit.Int32()
		category2     = gofakeit.Int32()
		manufacturer1 = model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.AppName(),
		}
		manufacturer2 = model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.AppName(),
		}
		tag1 = []string{gofakeit.Gamertag(), gofakeit.Gamertag()}
		tag2 = []string{gofakeit.Gamertag(), gofakeit.Gamertag()}

		uuids                 = []string{uuid1, uuid2}
		names                 = []string{name1, name2}
		categories            = []int32{category1, category2}
		manufacturerCountries = []string{manufacturer1.Country, manufacturer2.Country}
		tags                  = []string{tag1[0], tag2[0]}

		partsFilter = model.PartsFilter{
			Uuids:                 uuids,
			Names:                 names,
			Categories:            categories,
			ManufacturerCountries: manufacturerCountries,
			Tags:                  tags,
		}
		repoError = gofakeit.Error()
	)

	s.repository.On("ListParts", s.ctx, partsFilter).Return(nil, repoError)

	resp, err := s.service.ListParts(s.ctx, partsFilter)
	s.Require().Error(err)
	s.Require().Empty(resp)
	s.Require().Equal(repoError, err)
}
