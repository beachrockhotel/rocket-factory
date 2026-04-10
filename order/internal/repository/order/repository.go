package order

import (
	"sync"

	def "github.com/beachrockhotel/rocket-factory/order/internal/repository"
	"github.com/beachrockhotel/rocket-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*orderRepository)(nil)

type orderRepository struct {
	mtx  sync.RWMutex
	data map[string]model.OrderDTO
}

func NewRepository() *orderRepository {
	return &orderRepository{
		data: make(map[string]model.OrderDTO),
	}
}
