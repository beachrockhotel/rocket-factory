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

type OrderRequest struct {
	UserUUID  string
	PartUUIDs []string
}
