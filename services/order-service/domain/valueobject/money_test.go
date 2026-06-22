package valueobject_test

import (
	"errors"
	"testing"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	money "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
)

func TestMoney(t *testing.T) {
	t.Run("Tests for NewMoney", func(t *testing.T) {
		tests := []struct {
			name          string
			amount        int64
			currency      string
			expectedError error
		}{
			{"Valid Money", 100, "USD", nil},
			{"Negative Amount", -50, "USD", domainErrors.ErrNegativeAmount},
			{"Zero Amount", 0, "USD", domainErrors.ErrNegativeAmount},
			{"Invalid Currency", 100, "EUR", domainErrors.ErrInvalidCurrency},
			{"Empty Currency", 100, "", domainErrors.ErrInvalidCurrency},
			{"Valid BRL", 200, "BRL", nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m, err := money.NewMoney(tt.amount, tt.currency)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}
				if err == nil {
					if m.Amount() != tt.amount {
						t.Errorf("Expected amount: %d, got: %d", tt.amount, m.Amount())
					}
					if m.Currency() != tt.currency {
						t.Errorf("Expected currency: %s, got: %s", tt.currency, m.Currency())
					}
				}
			})
		}
	})

	t.Run("Test for NewMoneyMust", func(t *testing.T) {
		tests := []struct {
			name        string
			amount      int64
			currency    string
			expectPanic bool
		}{
			{"Valid Money", 100, "USD", false},
			{"Negative Amount", -50, "USD", true},
			{"Zero Amount", 0, "USD", true},
			{"Invalid Currency", 100, "EUR", true},
			{"Empty Currency", 100, "", true},
			{"Valid BRL", 200, "BRL", false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.expectPanic {
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("Expected panic but did not occur")
						}
					}()
				}
				m := money.NewMoneyMust(tt.amount, tt.currency)
				if !tt.expectPanic {
					if m.Amount() != tt.amount {
						t.Errorf("Expected amount: %d, got: %d", tt.amount, m.Amount())
					}
					if m.Currency() != tt.currency {
						t.Errorf("Expected currency: %s, got: %s", tt.currency, m.Currency())
					}
				}
			})
		}
	})

	t.Run("Test Money Amount and Currency methods", func(t *testing.T) {
		tests := []struct {
			name     string
			m        money.Money
			expected int64
			currency string
		}{
			{"USD Money", money.NewMoneyMust(100, "USD"), 100, "USD"},
			{"BRL Money", money.NewMoneyMust(200, "BRL"), 200, "BRL"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.m.Amount() != tt.expected {
					t.Errorf("Expected amount: %d, got: %d", tt.expected, tt.m.Amount())
				}
				if tt.m.Currency() != tt.currency {
					t.Errorf("Expected currency: %s, got: %s", tt.currency, tt.m.Currency())
				}
			})
		}
	})

	t.Run("Test Money Equals method", func(t *testing.T) {
		tests := []struct {
			name     string
			m1       money.Money
			m2       money.Money
			expected bool
		}{
			{"Equal Money", money.NewMoneyMust(100, "USD"), money.NewMoneyMust(100, "USD"), true},
			{"Different Amount", money.NewMoneyMust(100, "USD"), money.NewMoneyMust(200, "USD"), false},
			{"Different Currency", money.NewMoneyMust(100, "USD"), money.NewMoneyMust(100, "BRL"), false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.m1.Equals(tt.m2) != tt.expected {
					t.Errorf("Expected equality: %v, got: %v", tt.expected, tt.m1.Equals(tt.m2))
				}
			})
		}
	})

	t.Run("Test Money Multiply method", func(t *testing.T) {
		tests := []struct {
			name          string
			m             money.Money
			quantity      int64
			expectedMoney money.Money
			expectedError error
		}{
			{"Valid Multiply", money.NewMoneyMust(100, "USD"), 2, money.NewMoneyMust(200, "USD"), nil},
			{"Invalid Quantity", money.NewMoneyMust(100, "USD"), -1, money.Money{}, domainErrors.ErrInvalidQuantity},
			{"Zero Quantity", money.NewMoneyMust(100, "USD"), 0, money.Money{}, domainErrors.ErrInvalidQuantity},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := tt.m.Multiply(tt.quantity)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}
				if err == nil && !result.Equals(tt.expectedMoney) {
					t.Errorf("Expected money: %v, got: %v", tt.expectedMoney, result)
				}
			})
		}
	})

	t.Run("Test Money String method", func(t *testing.T) {
		tests := []struct {
			name     string
			m        money.Money
			expected string
		}{
			{"USD Money", money.NewMoneyMust(100, "USD"), "100 USD"},
			{"BRL Money", money.NewMoneyMust(200, "BRL"), "200 BRL"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.m.String() != tt.expected {
					t.Errorf("Expected string: %s, got: %s", tt.expected, tt.m.String())
				}
			})
		}
	})

	t.Run("Test for MarshalJSON", func(t *testing.T) {
		tests := []struct {
			name     string
			m        money.Money
			expected string
		}{
			{"USD Money", money.NewMoneyMust(100, "USD"), `{"amount":100,"currency":"USD"}`},
			{"BRL Money", money.NewMoneyMust(200, "BRL"), `{"amount":200,"currency":"BRL"}`},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				jsonData, err := tt.m.MarshalJSON()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if string(jsonData) != tt.expected {
					t.Errorf("Expected JSON: %s, got: %s", tt.expected, string(jsonData))
				}
			})
		}
	})

	t.Run("Test for UnmarshalJSON", func(t *testing.T) {
		tests := []struct {
			name          string
			jsonData      string
			expectedMoney money.Money
			expectedError error
		}{
			{"Valid JSON", `{"amount":100,"currency":"USD"}`, money.NewMoneyMust(100, "USD"), nil},
			{"Negative Amount", `{"amount":-50,"currency":"USD"}`, money.Money{}, domainErrors.ErrNegativeAmount},
			{"Invalid Currency", `{"amount":100,"currency":"EUR"}`, money.Money{}, domainErrors.ErrInvalidCurrency},
			{"Empty Currency", `{"amount":100,"currency":""}`, money.Money{}, domainErrors.ErrInvalidCurrency},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var m money.Money
				err := m.UnmarshalJSON([]byte(tt.jsonData))
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}
				if err == nil && !m.Equals(tt.expectedMoney) {
					t.Errorf("Expected money: %v, got: %v", tt.expectedMoney, m)
				}
			})
		}
	})
}
