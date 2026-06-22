package valueobject

import (
	"encoding/json"
	"fmt"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
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
		return Money{}, domainErrors.ErrNegativeAmount
	}
	if !validCurrencies[currency] {
		return Money{}, domainErrors.ErrInvalidCurrency
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func NewMoneyMust(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
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
		return Money{}, domainErrors.ErrInvalidQuantity
	}
	if m.amount <= 0 {
		return Money{}, domainErrors.ErrInvalidAmount
	}
	if !validCurrencies[m.currency] {
		return Money{}, domainErrors.ErrInvalidCurrency
	}
	return Money{
		amount:   m.amount * quantity,
		currency: m.currency,
	}, nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	auxMoney := struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}{
		Amount:   m.amount,
		Currency: m.currency,
	}
	return json.Marshal(auxMoney)
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
		return domainErrors.ErrNegativeAmount
	}
	if !validCurrencies[alias.Currency] {
		return domainErrors.ErrInvalidCurrency
	}

	m.amount = alias.Amount
	m.currency = alias.Currency
	return nil
}
