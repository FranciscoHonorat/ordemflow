package valueobject_test

import (
	"errors"
	"testing"

	"github.com/FranciscoHonorat/ordemflow/services/order-service/domain/domainErrors"
	address "github.com/FranciscoHonorat/ordemflow/services/order-service/domain/valueobject"
)

func TestAddress(t *testing.T) {
	t.Run("Tests for NewAddress", func(t *testing.T) {
		tests := []struct {
			name           string
			cep            string
			street         string
			neighborhood   string
			number         int64
			referencePoint string
			complement     string
			expectedError  error
		}{
			{"valid Address", "12345678", "Street", "Neighborhood", 02, "reference", "complement", nil},
			{"Invalid CEP", "123456", "Street", "Neighborhood", 02, "reference", "complement", domainErrors.ErrInvalidCEP},
			{"Invalid Street", "12345678", "", "Neighborhood", 02, "reference", "complement", domainErrors.ErrFieldEmpty},
			{"Invalid Neighborhood", "12345678", "Street", "", 02, "reference", "complement", domainErrors.ErrFieldEmpty},
			{"Invalid Number", "12345678", "Street", "Neighborhood", -1, "reference", "complement", domainErrors.ErrInvalidNumber},
			{"Invalid ReferencePoint", "12345678", "Street", "Neighborhood", 02, "", "complement", domainErrors.ErrFieldEmpty},
			{"Invalid Complement", "12345678", "Street", "Neighborhood", 02, "reference", "", domainErrors.ErrFieldEmpty},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, err := address.NewAddress(tt.cep, tt.street, tt.neighborhood, tt.number, tt.referencePoint, tt.complement)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				}

				if err == nil {
					if a.Cep() != tt.cep {
						t.Errorf("Expected CEP: %s, got: %s", tt.cep, a.Cep())
					}
					if a.Street() != tt.street {
						t.Errorf("Expected Street: %s, got: %s", tt.street, a.Street())
					}
					if a.Neighborhood() != tt.neighborhood {
						t.Errorf("Expected Neighborhood: %s, got: %s", tt.neighborhood, a.Neighborhood())
					}
					if a.Number() != tt.number {
						t.Errorf("Expected Number: %d, got: %d", tt.number, a.Number())
					}
					if a.ReferencePoint() != tt.referencePoint {
						t.Errorf("Expected Reference Point: %s, got: %s", tt.referencePoint, a.ReferencePoint())
					}
					if a.Complement() != tt.complement {
						t.Errorf("Expected Complement: %s, got: %s", tt.complement, a.Complement())
					}
				}
			})
		}
	})

	t.Run("Test for MarshalJSON", func(t *testing.T) {
		tests := []struct {
			name     string
			a        address.Address
			expected string
		}{
			{
				name:     "valid Address",
				a:        address.NewAddressMust("12345678", "Street", "Neighborhood", 2, "reference", "complement"),
				expected: `{"cep":"12345678","street":"Street","neighborhood":"Neighborhood","number":2,"referencePoint":"reference","complement":"complement"}`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				jsonData, err := tt.a.MarshalJSON()
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
			name            string
			jsonData        string
			expectedAddress address.Address
			expectedError   error
		}{
			{"Valid JSON", `{"cep":"12345678","street":"Street","neighborhood":"Neighborhood","number":2,"referencePoint":"reference","complement":"complement"}`, address.NewAddressMust("12345678", "Street", "Neighborhood", 2, "reference", "complement"), nil},
			{"Invalid Street", `{"cep":"12345678","street":"","neighborhood":"Neighborhood","number":2,"referencePoint":"reference","complement":"complement"}`, address.Address{}, domainErrors.ErrFieldEmpty},
			{"Invalid Neighborhood", `{"cep":"12345678","street":"Street","neighborhood":"","number":2,"referencePoint":"reference","complement":"complement"}`, address.Address{}, domainErrors.ErrFieldEmpty},
			{"Invalid Number", `{"cep":"12345678","street":"Street","neighborhood":"Neighborhood","number":-1,"referencePoint":"reference","complement":"complement"}`, address.Address{}, domainErrors.ErrInvalidNumber},
			{"Invalid ReferencePoint", `{"cep":"12345678","street":"Street","neighborhood":"Neighborhood","number":2,"referencePoint":"","complement":"complement"}`, address.Address{}, domainErrors.ErrFieldEmpty},
			{"Invalid Complement", `{"cep":"12345678","street":"Street","neighborhood":"Neighborhood","number":2,"referencePoint":"reference","complement":""}`, address.Address{}, domainErrors.ErrFieldEmpty},
		}
		for _, tt := range tests {
			var a address.Address
			err := a.UnmarshalJSON([]byte(tt.jsonData))
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}
		}
	})
}
