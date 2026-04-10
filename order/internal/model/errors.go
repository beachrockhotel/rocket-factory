package model

import "github.com/go-faster/errors"

var (
	ErrOrderTimeout        = errors.New("order service timeout")
	ErrInternalServerError = errors.New("order service errors")
	ErrBadRequest          = errors.New("some parts not found")
	ErrNotFound            = errors.New("order not found")
	ErrConflict            = errors.New("order cannot be paid")
)
