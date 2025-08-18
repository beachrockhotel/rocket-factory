package main

import (
	"math"

	inventoryV1 "github.com/beachrockhotel/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *inventoryService) initParts() {
	parts := generateParts()
	for _, part := range parts {
		s.inv[part.Uuid] = part
	}
}

func generateParts() []*inventoryV1.Part {
	names := []string{
		"Main Engine", "Reserve Engine", "Thruster", "Fuel Tank",
		"Left Wing", "Right Wing", "Window A", "Window B",
		"Control Module", "Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit", "Backup propulsion unit", "Thruster for fine adjustments",
		"Main fuel tank", "Left aerodynamic wing", "Right aerodynamic wing",
		"Front viewing window", "Side viewing window", "Flight control module",
		"Stabilization fin",
	}

	var parts []*inventoryV1.Part
	for i := 0; i < gofakeit.Number(10, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryV1.Category(gofakeit.Number(1, 4)), //nolint:gosec
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}
	return parts
}

func generateDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 5); i++ {
		tags = append(tags, gofakeit.HackerNoun()) // вместо эмодзи
	}
	return tags
}

func generateMetadata() map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value)
	for i := 0; i < gofakeit.Number(1, 5); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}
	return metadata
}

func generateMetadataValue() *inventoryV1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: gofakeit.Word()}}
	case 1:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: int64(gofakeit.Number(1, 100))}}
	case 2:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: roundTo(gofakeit.Float64Range(1, 100))}}
	case 3:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: gofakeit.Bool()}}
	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
