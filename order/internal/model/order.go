package model

import "time"

type OrderDTO struct {
	UUID            string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type CreateOrderRequest struct {
	UserUUID  string
	PartUUIDs []string
}

type CreateOrderParams struct {
	UserUUID   string
	PartUUIDs  []string
	TotalPrice float64
	Status     string
}

const (
	OrderStatusPendingPayment = "PENDING_PAYMENT"
	OrderStatusPaid           = "PAID"
	OrderStatusCancelled      = "CANCELLED"
)
