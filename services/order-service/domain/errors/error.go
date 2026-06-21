package errors

import "errors"

var (
	ErrNegativeAmount    = errors.New("amount cannot be negative")
	ErrInvalidCurrency   = errors.New("invalid currency")
	ErrInvalidQuantity   = errors.New("quantity invalid")
	ErrInvalidCEP        = errors.New("The postal code is exactly 8 digits.")
	ErrFieldEmpty        = errors.New("This field cannot be empty.")
	ErrInvalidNumber     = errors.New("The number needs to be greater than 0.")
	ErrInvalidID         = errors.New("ID cannot be empty.")
	ErrRecalculeTotal    = errors.New("Error recalculating total price.")
	ErrInvalidStatus     = errors.New("Invalid order status.")
	ErrInvalidOrderID    = errors.New("Invalid order ID.")
	ErrInvalidProductID  = errors.New("Invalid product ID.")
	ErrInvalidCustomerID = errors.New("Invalid customer ID.")
	ErrInvalidAmount     = errors.New("Amount must be greater than 0.")
)
