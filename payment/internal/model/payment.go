package model

type PayOrderRequest struct {
	OrderUuid     string
	UserUUID      string
	PaymentMethod string
}
