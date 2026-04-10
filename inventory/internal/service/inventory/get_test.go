package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/beachrockhotel/rocket-factory/inventory/internal/model"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		uuid          = gofakeit.UUID()
		name          = gofakeit.Name()
		description   = gofakeit.ProductDescription()
		price         = gofakeit.Float64()
		stockQuantity = gofakeit.Int64()
		category      = gofakeit.Int32()
		dimension     = model.Dimensions{
			Length: gofakeit.Float64(),
			Width:  gofakeit.Float64(),
			Height: gofakeit.Float64(),
			Weight: gofakeit.Float64(),
		}
		manufacturer = model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.AppName(),
		}
		Tags     = []string{gofakeit.Gamertag(), gofakeit.Gamertag()}
		Metadata = map[string]any{
			gofakeit.Name(): gofakeit.Country(),
			gofakeit.Name(): gofakeit.Country(),
			gofakeit.Name(): gofakeit.Country(),
		}
		CreatedAt = gofakeit.Date()
		UpdatedAt = gofakeit.Date()
		part      = model.Part{
			Uuid:          uuid,
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      category,
			Dimension:     dimension,
			Manufacturer:  manufacturer,
			Tags:          Tags,
			Metadata:      Metadata,
			CreatedAt:     CreatedAt,
			UpdatedAt:     &UpdatedAt,
		}
	)
	s.repository.On("GetPart", s.ctx, uuid).Return(part, nil)

	resp, err := s.service.GetPart(s.ctx, uuid)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)
	s.Require().Equal(part, resp)
}

func (s *ServiceSuite) TestGetPartError() {
	var (
		repoError = gofakeit.Error()
		uuid      = gofakeit.UUID()
	)
	s.repository.On("GetPart", s.ctx, uuid).Return(model.Part{}, repoError)

	resp, err := s.service.GetPart(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().Empty(resp)
	s.Require().Equal(repoError, err)
}
