package model

type PayOrderRequest struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod string
}
