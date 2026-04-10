package inventory

import (
	"sync"

	def "github.com/beachrockhotel/rocket-factory/inventory/internal/repository"
	repoModel "github.com/beachrockhotel/rocket-factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mtx   sync.RWMutex
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	r := &repository{
		parts: make(map[string]repoModel.Part),
	}
	r.initParts()

	return r
}
