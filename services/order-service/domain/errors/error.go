package errors

import "errors"

var (
	ErrNegativeAmount  = errors.New("amount cannot be negative")
	ErrInvalidCurrency = errors.New("invalid currency")
)
