package valueobject

import (
	"encoding/json"
	"fmt"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/errors"
)

var validCurrencies = map[string]bool{
	"BRL": true,
	"USD": true,
}

type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount <= 0 {
		return Money{}, errors.ErrNegativeAmount
	}
	if !validCurrencies[currency] {
		return Money{}, errors.ErrInvalidCurrency
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount, m.currency)
}

func (m Money) Multiply(quantity int64) (Money, error) {
	if quantity <= 0 {
		return Money{}, errors.ErrInvalidQuantity
	}
	if m.amount <= 0 {
		return Money{}, errors.ErrInvalidAmount
	}
	if !validCurrencies[m.currency] {
		return Money{}, errors.ErrInvalidCurrency
	}
	return Money{
		amount:   m.amount * quantity,
		currency: m.currency,
	}, nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"amount": %d, "currency": "%s"}`, m.amount, m.currency)), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var alias struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.Amount <= 0 {
		return errors.ErrNegativeAmount
	}
	if !validCurrencies[alias.Currency] {
		return errors.ErrInvalidCurrency
	}

	m.amount = alias.Amount
	m.currency = alias.Currency
	return nil
}
