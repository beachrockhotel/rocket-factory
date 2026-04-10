package model

import "time"

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int32
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]any
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

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

type PartsFilter struct {
	Uuids                 *[]string
	Names                 *[]string
	Categories            *[]int32
	ManufacturerCountries *[]string
	Tags                  *[]string
}
