package repository

import (
	"sync"
)

type paymentRepository struct {
	mtx sync.Mutex
}

func NewPaymentRepository() *paymentRepository {
	return &paymentRepository{}
}
