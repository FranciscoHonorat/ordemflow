package errors

import "errors"

var (
	ErrNegativeAmount  = errors.New("amount cannot be negative")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidQuantity = errors.New("quantity invalid")
	ErrInvalidCEP      = errors.New("The postal code is exactly 8 digits.")
	ErrFieldEmpty      = errors.New("This field cannot be empty.")
	ErrInvalidNumber   = errors.New("The number needs to be greater than 0.")
)
