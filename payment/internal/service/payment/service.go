package payment

import (
	"github.com/beachrockhotel/rocket-factory/payment/internal/repository"
)

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *paymentService {
	return &paymentService{
		repo: repo,
	}
}
